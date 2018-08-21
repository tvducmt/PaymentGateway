package models

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

// Receipt ...
type Receipt struct {
	ID          int       `json:"id" db:"id"`
	RawReceipt  string    `json:"raw_receipt" db:"raw_receipt"`
	Amount      float64   `json:"parsed_amount" db:"parsed_amount"`
	Account     string    `json:"parsed_account" db:"parsed_account"`
	PhoneNumber string    `json:"phonenumber" db:"phonenumber"`
	CreateAt    time.Time `json:"create_at" db:"create_at"`
}

// NewReceipt ...
func NewReceipt(rawReceipt, phonenumber string, createAt time.Time) *Receipt {
	return &Receipt{
		RawReceipt:  rawReceipt,
		PhoneNumber: phonenumber,
		CreateAt:    createAt,
	}
}

// Add ...
func (r *Receipt) Add(tx *sql.Tx) (int, error) {
	var receiptID int
	query := `INSERT INTO receipt (raw_receipt, phonenumber, create_at) VALUES ($1, $2, $3) RETURNING id`
	err := tx.QueryRow(query, r.RawReceipt, r.PhoneNumber, r.CreateAt).Scan(&receiptID)
	if err != nil {
		return 0, err
	}
	return receiptID, nil
}

// AddWithCoupon ...
func (r *Receipt) AddWithCoupon(tx *sql.Tx, value float64) (int, error) {
	var receiptID int
	query := `INSERT INTO receipt (coupon_value) VALUES ($1) RETURNING id`
	err := tx.QueryRow(query, -value).Scan(&receiptID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return receiptID, nil
}

// Covert a string to float
// Chuyển đổi kiểu string sang kiểu float
func covertStringToFloat(amount string) float64 {
	amountString := strings.Replace(amount, ",", "", -1)
	amountMoney, err := strconv.Atoi(amountString)
	if err != nil {
		return 0
	}
	amountFloat := float64(amountMoney)
	return amountFloat
}

// UpdateReceiptByreceiptID ...
//khi nhận receipt từ client thì hệ thống kiểm tra receipt có hợp lên không
// Pased các thông tin và update vào bảng receipt
func UpdateReceiptByreceiptID(tx *sql.Tx, receiptID int, dataPasred map[string]string) error {
	amount := covertStringToFloat(dataPasred["amount"])
	balance := covertStringToFloat(dataPasred["balance"])
	if dataPasred["second"] == "" {
		timeReceipt := dataPasred["year"] + "-" + dataPasred["month"] + "-" + dataPasred["day"] + "T" + dataPasred["hour"] + ":" + dataPasred["minute"]
		timestamp, err := time.Parse("2006-01-02T15:04", timeReceipt)
		if err != nil {
			return err
		}
		queryReceipt := `UPDATE receipt
		SET   parsed_amount=$1, parsed_account=$2, parsed_code=$3, parsed_balance=$4, sys_create_at=$5
		WHERE id= $6;`
		_, err = tx.Exec(queryReceipt, amount, dataPasred["bank_number"], dataPasred["transaction_code"], balance, timestamp, receiptID)
		if err != nil {
			return err
		}

	} else {
		timeReceipt := dataPasred["year"] + "-" + dataPasred["month"] + "-" + dataPasred["day"] + "T" + dataPasred["hour"] + ":" + dataPasred["minute"] + ":" + dataPasred["second"]
		timestamp, err := time.Parse("2006-01-02T15:04:05", timeReceipt)
		if err != nil {
			return err
		}
		query := `UPDATE receipt
		SET   parsed_amount=$1, parsed_account=$2, parsed_code=$3, parsed_balance=$4, sys_create_at=$5
		WHERE id= $6;`
		_, err = tx.Exec(query, amount, dataPasred["bank_number"], dataPasred["transaction_code"], balance, timestamp, receiptID)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateReceiptByID ...
func UpdateReceiptByID(tx *sql.Tx, id int) (float64, error) {
	var amount sql.NullFloat64
	querySelect := "select parsed_amount from receipt where id = $1"
	err := tx.QueryRow(querySelect, id).Scan(&amount)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	queryUpdate := "update receipt set coupon_value = $1 where id = $2"
	_, err = tx.Exec(queryUpdate, -amount.Float64, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return amount.Float64, nil
}
