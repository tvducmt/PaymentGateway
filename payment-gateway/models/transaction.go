package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	lock "github.com/bsm/redis-lock"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Transaction represents a transaction function found in a source files
type Transaction struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	ReceiptID int       `json:"receipt_id" db:"receipt_id"`
	MethodID  int       `json:"method_id" db:"method_id"`
	GatewayID int       `json:"gateway_id" db:"gateway_id"`
	Status    string    `json:"status" db:"status"`
	Code      string    `json:"code" db:"code"`
	ChargeID  string    `json:"charge_id" db:"charge_id"`
	CreateAt  time.Time `json:"create_at" db:"create_at"`
}

// NewTransaction this function to create transaction
// Return transaction
func NewTransaction(userID, methodID, gatewayID int) *Transaction {
	return &Transaction{
		UserID:    userID,
		MethodID:  methodID,
		GatewayID: gatewayID,
		Status:    "pending",
	}
}

// NewCreditCard ..
func NewCreditCard(userID, methodID int, chargeID string) *Transaction {
	return &Transaction{
		UserID:   userID,
		MethodID: methodID,
		Status:   "failed",
		ChargeID: chargeID,
	}
}

// FindTransactionByID this function to find transaction  by id
func FindTransactionByID(db *sqlx.DB, txID int) (*Transaction, error) {
	var tx Transaction
	err := db.Get(&tx, "SELECT * FROM transaction WHERE id = $1")
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// swagger:model  TransactionUser
type TransactionUser struct {
	ID          int             `json:"id" db:"id"`
	CreateAt    time.Time       `json:"create_at" db:"create_at"`
	Status      string          `json:"status" db:"status"`
	MethodName  sql.NullString  `json:"method_name" db:"method_name"`
	RawReceipt  sql.NullString  `json:"raw_receipt" db:"raw_receipt"`
	Code        sql.NullString  `json:"code" db:"code"`
	Amount      sql.NullFloat64 `json:"parsed_amount" db:"parsed_amount"`
	Account     sql.NullString  `json:"parsed_account" db:"parsed_account"`
	PhoneNumber sql.NullString  `json:"phonenumber" db:"phonenumber"`
	ReceiptDate pq.NullTime     `json:"sys_create_at" db:"sys_create_at"`
	CStatus     sql.NullString  `json:"coupon_status" db:"coupon_status"`
	CValue      sql.NullFloat64 `json:"coupon_value" db:"coupon_value"`
	CCode       sql.NullString  `json:"coupon_code" db:"coupon_code"`
	CCurrency   sql.NullString  `json:"coupon_currency" db:"coupon_currency"`
	ChargeID    sql.NullString  `json:"charge_id" db:"charge_id"`
	CCreateAt   pq.NullTime     `json:"coupon_create_at" db:"coupon_create_at"`
}

// FindTransactionByUserID this function to find transaction  by userid
func FindTransactionByUserID(db *sqlx.DB, userID int) ([]TransactionUser, error) {
	var txs []TransactionUser
	query := `SELECT tx.id id, 
				tx.create_at create_at, 
				tx.status status,
				pm.name method_name,
				r.raw_receipt raw_receipt, 
				tx.code code,
				r.parsed_amount parsed_amount,
				r.parsed_account parsed_account,
				r.phonenumber phonenumber,
				r.sys_create_at sys_create_at,
				c.status coupon_status,
				c.value coupon_value,
				c.code coupon_code,
				c.currency coupon_currency,
				tx.charge_id,
				c.create_at coupon_create_at
				FROM transaction tx 
				LEFT JOIN receipt r ON tx.receipt_id = r.id 
				LEFT JOIN coupon c ON tx.id = c.tx_id
				LEFT JOIN payment_method pm ON tx.method_id = pm.id
				WHERE tx.user_id = $1;`
	err := db.Select(&txs, query, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return txs, nil
}

// FindFullTransactionByID ...
func FindFullTransactionByID(db *sqlx.DB, txID int) ([]TransactionUser, error) {
	var txs []TransactionUser
	query := `SELECT tx.id id, 
				tx.create_at create_at, 
				tx.status status,
				pm.name method_name,
				r.raw_receipt raw_receipt, 
				tx.code code,
				r.parsed_amount parsed_amount,
				r.parsed_account parsed_account,
				r.phonenumber phonenumber,
				r.sys_create_at sys_create_at,
				c.status coupon_status,
				c.value coupon_value,
				c.code coupon_code,
				c.currency coupon_currency,
				tx.charge_id,
				c.create_at coupon_create_at
				FROM transaction tx 
				LEFT JOIN receipt r ON tx.receipt_id = r.id 
				LEFT JOIN coupon c ON tx.id = c.tx_id
				LEFT JOIN payment_method pm ON tx.method_id = pm.id
				WHERE tx.id = $1;`
	err := db.Select(&txs, query, txID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return txs, nil
}

type Money struct {
	TotalAmount sql.NullFloat64
	Currency    sql.NullString
}

// SumBalanceByUserID total balance
func SumBalanceByUserID(db *sqlx.DB, userID int) ([]Money, error) {
	money := Money{}
	moneys := []Money{}
	query := `SELECT sum(c.value) totalAmount, c.currency currency
			FROM  transaction tx LEFT JOIN coupon c ON c.tx_id = tx.id WHERE tx.status = 'confirmed' AND tx.user_id = $1 
			GROUP BY currency`
	rows, err := db.Query(query, userID)
	if err != nil {
		fmt.Println("Error", err)
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&money.TotalAmount, &money.Currency)
		if err != nil {
			fmt.Println("Error scan", err)
		}
		moneys = append(moneys, money)
	}
	return moneys, nil
}

var datalock = "key"

func getTxCode(db *sqlx.DB) (string, error) {
	var txCode string
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
	})
	defer client.Close()
	for {
		lock, err := lock.Obtain(client, datalock, &lock.Options{LockTimeout: 10 * time.Second})
		if err != nil {
			continue
		} else if lock == nil {
			fmt.Println("ERROR: could not obtain lock")
		}
		querySelect := "select code from codetable where active=false and id < $1"
		err = db.QueryRow(querySelect, CodeQuantity+1).Scan(&txCode)
		if err != nil {
			fmt.Println("Error ", err)
			return "", err
		}
		queryUpdate := "UPDATE codetable SET active = true WHERE code = $1"
		_, err = db.Exec(queryUpdate, txCode)
		if err != nil {
			fmt.Println("Error ", err)
			return "", err
		}
		err = lock.Unlock()
		if err != nil {
			fmt.Println("Error ", err)
		}
		return txCode, nil
	}
}

// Add transaction when have requests
func (t *Transaction) Add(db *sqlx.DB) (string, error) {
	txcode, err := getTxCode(db)
	if err != nil {
		return "", err
	}
	_, err = db.Exec(`INSERT INTO transaction (user_id, method_id, gateway_id, code, status) VALUES ($1, $2, $3, $4, $5)`, t.UserID, t.MethodID, t.GatewayID, txcode, t.Status)
	if err != nil {
		return "", err
	}
	return txcode, nil
}

// AddCreditCard transaction when have requests
func (t *Transaction) AddCreditCard(db *sqlx.DB) error {
	err := db.QueryRow(`INSERT INTO transaction (user_id, method_id, status, charge_id) VALUES ($1, $2, $3, $4) RETURNING id`, t.UserID, t.MethodID, t.Status, t.ChargeID).Scan(&t.ID)
	if err != nil {
		return err
	}
	return nil
}

// AddWithCoupon transaction when have requests
func (t *Transaction) AddWithCoupon(tx *sql.Tx, userID int) (int, error) {
	var txID int
	err := tx.QueryRow(`INSERT INTO transaction (user_id) VALUES ($1) RETURNING id`, userID).Scan(&txID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return txID, nil
}

//UpdateTXWithCoupon ..
func (t *Transaction) UpdateTXWithCoupon(tx *sql.Tx, txID int, receiptID int) error {
	queryUpdate := `UPDATE transaction
		SET  status='confirmed', receipt_id=$1
		WHERE id = $2 AND status = 'pending';`
	_, err := tx.Exec(queryUpdate, receiptID, txID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// UpdateTXByCode ...
func (t *Transaction) UpdateTXByCode(tx *sql.Tx, receiptID int, codeID string) (int, error) {
	var count int
	row := tx.QueryRow(`select count(code) from transaction where code= $1 and method_id=1`, codeID)
	err := row.Scan(&count)
	if count == 0 {
		return 0, errors.New("CodeID no match any transaction")
	}
	var txID int
	queryUpdate := `UPDATE transaction
		SET  status='confirmed', receipt_id=$1
		WHERE code = $2 AND status = 'pending' RETURNING id;`
	err = tx.QueryRow(queryUpdate, receiptID, codeID).Scan(&txID)
	if err != nil {
		return 0, err
	}
	return txID, nil
}

// UpdateTXByCode ...
func (t *Transaction) UpdateTXByContractID(tx *sql.Tx, contractID int, codeID string) (int, error) {
	var count int
	row := tx.QueryRow(`select count(code) from transaction where code= $1 and method_id=4`, codeID)
	err := row.Scan(&count)
	if count == 0 {
		return 0, errors.New("CodeID no match any transaction")
	}
	var txID int
	queryUpdate := `UPDATE transaction
		SET  status='confirmed', contract_id=$1
		WHERE code = $2 AND status = 'pending' RETURNING id;`
	err = tx.QueryRow(queryUpdate, contractID, codeID).Scan(&txID)
	if err != nil {
		return 0, err
	}
	return txID, nil
}

// UpdateTimeOut ....
func (t *Transaction) UpdateTimeOut(db *sqlx.DB, txID int) error {
	query := "UPDATE transaction SET status = 'timeout' WHERE id = $1"
	_, err := db.Exec(query, txID)
	if err != nil {
		return errors.New("Can't update timeout")
	}
	return nil
}

// UpdateSuccess ....
func UpdateSuccess(db *sqlx.DB, txID int) error {
	query := "UPDATE transaction SET status = 'successed' WHERE id = $1"
	_, err := db.Exec(query, txID)
	if err != nil {
		return errors.New("Don't update successed")
	}
	return nil
}

//CheckTransactionTimeOut ...
func (t *Transaction) CheckTransactionTimeOut(db *sqlx.DB) error {

	type tx struct {
		ID       int       ` db:"id"`
		CreateAt time.Time ` db:"create_at"`
	}
	txs := []tx{}
	query := "SELECT id, create_at FROM transaction WHERE status = 'pending'"
	err := db.Select(&txs, query)
	if err != nil {
		log.Println(err)
		return err
	}
	timeUnix := time.Now().Unix()
	//fmt.Println(nanos)
	for _, v := range txs {
		if (timeUnix - v.CreateAt.Unix()) > 3600 {
			t.UpdateTimeOut(db, v.ID)
		}

	}
	return nil
}

// FindTxIDByUserID ..
func (t *Transaction) FindTxIDByUserID(db *sqlx.DB, userID int) ([]int, error) {
	var listTxID []int
	query := "SELECT id FROM Transaction WHERE user_id = $1 AND (status = 'confirmed' OR status = 'successed')"
	err := db.Select(&listTxID, query, userID)
	if err != nil {
		return nil, err
	}
	return listTxID, nil
}
