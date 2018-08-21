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

// swagger:model Coupon
type Coupon struct {
	// ID định danh cho coupon
	//
	// min: 1
	ID int `json:"id"`

	// ID khóa ngoại đến model transaction
	//
	// min: 1
	TxID int `json:"tx_id"`

	// Trạng thái Coupon (spend, unspend)
	Status string `json:"status"`

	// Mã Coupon Code
	Code string `json:"code"`

	// Thời gian tạo Coupon
	CreateAt time.Time `json:"create_at"`

	// Giá trị của Coupon
	// min: 1
	Value float64 `json:"value"`

	// Đơn vị tiền tệ của Coupon. Mặc định là vnd
	Currency string `json:"currency"`
}

func getCouponCode(tx *sql.Tx) (string, error) {
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
		querySelect := "SELECT code FROM codetable WHERE active=false and id > $1 AND id < $2"
		err = tx.QueryRow(querySelect, CodeQuantity, 2*CodeQuantity+1).Scan(&txCode)
		if err != nil {
			tx.Rollback()
			fmt.Println("Error ", err)
			return "", err
		}
		queryUpdate := "UPDATE codetable SET active = true WHERE code = $1"
		_, err = tx.Exec(queryUpdate, txCode)
		if err != nil {
			tx.Rollback()
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

//NewCoupon ...
func NewCoupon(tx *sql.Tx, txID int, value float64, currency string) (*Coupon, error) {
	couponCode, err := getCouponCode(tx)
	if err != nil {
		tx.Rollback()
		fmt.Println("Model NewCoupon > Error ", err)
	}
	// ckeck exist coupon (newcoupon, select couponcode from coupon where )
	return &Coupon{
		TxID:     txID,
		Status:   "unspend",
		Code:     couponCode,
		Value:    value,
		Currency: currency,
	}, nil
}

// Add ...
func (c *Coupon) Add(tx *sql.Tx) (*Coupon, error) {

	var coupon = Coupon{}
	query := `INSERT INTO coupon (tx_id, code, value, currency) VALUES ($1, $2, $3, $4) RETURNING id, tx_id, status, code, create_at, value`
	err := tx.QueryRow(query, c.TxID, c.Code, c.Value, c.Currency).Scan(&coupon.ID, &coupon.TxID, &coupon.Status, &coupon.Code, &coupon.CreateAt, &coupon.Value)
	if err != nil {
		tx.Rollback()
		pqErr, ok := err.(*pq.Error)
		if ok {
			if pqErr.Code.Name() == "unique_violation" {
				return nil, errors.New("Không hợp lệ. Giao dịch này đã được sử dụng để tạo Coupon trước đây.")
			}
		}
		return nil, err
	}
	return &coupon, nil
}

// CreateCouponWithNewReceipt ...
func (c *Coupon) CreateCouponWithNewReceipt(tx *sql.Tx, userID int, txID int, value float64) (*Coupon, error) {
	var checkUser int //() Check txID có phải là của user đó không
	err := tx.QueryRow("select count(*) from transaction where id = $1 AND user_id = $2", txID, userID).Scan(&checkUser)
	if err != nil {
		return nil, err
	}
	if checkUser != 1 {
		tx.Rollback()
		return nil, errors.New("Transaction not match with user, Please try again")
	}
	coupon, err := NewCoupon(tx, txID, value, "vnd")
	if err != nil {
		log.Println("Create coupon failed", err)
		return nil, err
	}
	cp, err := coupon.Add(tx)
	if err != nil {
		log.Println("Add failed", err)
		tx.Rollback()
		return nil, err
	}

	return cp, nil

}

// FindCouponsByTxID ..
func (c *Coupon) FindCouponsByTxID(db *sqlx.DB, listTxID []int) ([]Coupon, error) {
	var listCoupon []Coupon
	query := `SELECT * FROM coupon WHERE tx_id = $1`
	for _, v := range listTxID {
		cp := Coupon{}
		err := db.QueryRow(query, v).Scan(&cp.ID, &cp.TxID, &cp.Status, &cp.Code, &cp.CreateAt, &cp.Value, &cp.Currency)
		if err != nil {
			continue
		} else {
			listCoupon = append(listCoupon, cp)
		}
	}
	return listCoupon, nil
}
