package main

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"gitlab.com/rockship/payment-gateway/models"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func EthereumMonitoring(db *sqlx.DB) {
	// TokenABI is the input ABI used to generate the binding from.
	const TokenABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"wallet\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_wallet\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"NewTransaction\",\"type\":\"event\"}]"
	domainRinkeby := os.Getenv("Rinkeby_DSN")
	client, err := ethclient.Dial(domainRinkeby)

	if err != nil {
		fmt.Println(err)
		return
	}
	contractAddress := common.HexToAddress(os.Getenv("ContractAddr"))
	query := ethereum.FilterQuery{}
	query.Addresses = make([]common.Address, 0)
	query.Addresses = append(query.Addresses, contractAddress)

	ctx := context.Background()
	var ch = make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(ctx, query, ch)

	abi, err := abi.JSON(strings.NewReader(string(TokenABI)))
	if err != nil {
		fmt.Println("Invalid abi:", err)
	}
	fmt.Println("Listening event ...")
	for {
		select {
		case err := <-sub.Err():
			log.Println(err)
		case eventLog := <-ch:
			txHash := common.HexToHash(eventLog.TxHash.String())
			tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
			if err != nil {
				log.Println(err)
			}
			if !isPending {
				var transferEvent struct {
					Sender common.Address
					Value  *big.Int
				}
				err = abi.Unpack(&transferEvent, "NewTransaction", eventLog.Data)
				if err != nil {
					log.Println("Failed to unpack", err)
					continue
				}
				fbalance := new(big.Float)
				fbalance.SetString(transferEvent.Value.String())
				ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
				etherString := ethValue.String()
				etherFloat, err := strconv.ParseFloat(etherString, 64)
				if err != nil {
					fmt.Println("Error ", err)
				}
				txEthereum := models.NewTxEthereum(transferEvent.Sender.Hex(), txHash.String(), etherFloat)
				id, err := txEthereum.Add(db)
				if id == 0 {
					fmt.Println("Error Add ", err)
					return
				}
				data := tx.Data()
				if data != nil {
					txID, err := transactionProccess(db, string(tx.Data()), id)
					if err != nil {
						fmt.Println("Error ", err)
					}
					txHandle, _ := db.Begin()
					cp, err := models.NewCoupon(txHandle, txID, txEthereum.Value, "eth")
					if err != nil {
						fmt.Println("Create coupon failed", err)
						//	w.WriteHeader(http.StatusInternalServerError)
					}
					coupon, err := cp.Add(txHandle)
					if err != nil {
						fmt.Println("Add failed", err)
						return
					}
					fmt.Println(coupon)
					txHandle.Commit()
				}
				fmt.Println("Successfull")

			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func transactionProccess(db *sqlx.DB, data string, contractID int) (int, error) {
	rules, err := models.LoadRegex(db)
	if err != nil {
		return 0, err
	}
	checkRule, err := models.CheckRule(data, rules)
	if err != nil {
		return 0, err
	}
	txHandle, err := db.Begin()
	var txID int
	if checkRule != nil {
		txcode := checkRule["transaction_code"]
		tx := models.Transaction{}
		txID, err = tx.UpdateTXByContractID(txHandle, contractID, txcode)
		if err != nil {
			fmt.Println("Error UpdateTXContract", err)
			return 0, err
		}
	}
	txHandle.Commit()
	return txID, nil
}
