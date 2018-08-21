package models

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	stripe "github.com/stripe/stripe-go"
)

// Customer ...
type Customer struct {
	ID          int    `json:"id" db:"id"`
	UserID      int    `json:"user_id" db:"user_id"`
	CusStripeID string `json:"cus_stripe_id" db:"cus_stripe_id"`
	Fingerprint string `json:"fingerprint" db:"fingerprint"`
}

// swagger:model  CustomerUser
type CustomerUser struct {
	// ID định danh cho customer
	//
	// min: 1
	ID int `json:"id"`

	// UserID khóa ngoại đến model user
	//
	// min: 1
	UserID int `json:"user_id"`

	// Bốn số cuối thẻ
	Last4 string `json:"last4"`

	// Tháng hết hạn
	ExpMonth uint8 `json:"exp_month"`

	// Năm hết hạn
	ExpYear uint16 `json:"exp_year"`

	// Mã bưu điện
	AddressZip string `json:"address_zip"`

	// Loại thẻ (Visa, MasterCard, ...)
	Brand stripe.CardBrand `json:"brand"`

	// Hình thức (Credit, Debit, ...)
	Funding stripe.CardFunding `json:"funding"`

	// Chữ ký trên thẻ
	// Unique: true
	Fingerprint string `json:"fingerprint"`
}

// NewCustomer ...
func NewCustomer(userID int, cusStripeID string) *Customer {
	var cus Customer
	cus.UserID = userID
	cus.CusStripeID = cusStripeID
	return &cus
}

// AddCustomerToUser ... add customer_id from stripe.js
func (c *Customer) AddCustomerToUser(db *sqlx.DB) (int, error) {
	var id int
	err := db.QueryRow(`INSERT INTO customer (user_id, cus_stripe_id, fingerprint) VALUES($1, $2, $3) RETURNING id`, c.UserID, c.CusStripeID, c.Fingerprint).Scan(&id)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			if pqErr.Code.Name() == "unique_violation" {
				return -1, nil
			}
		}
		return 0, err
	}
	return id, nil
}

// FindCustomerByUserID ...
func FindCustomerByUserID(db *sqlx.DB, userID int) ([]Customer, error) {
	var cus []Customer
	query := `SELECT id, user_id, cus_stripe_id FROM customer WHERE user_id = $1`
	err := db.Select(&cus, query, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return cus, nil
}

// FindCustomerByID ...
func FindCustomerByID(db *sqlx.DB, id int) (*Customer, error) {
	var cus Customer
	err := db.Get(&cus, `SELECT * FROM customer WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &cus, nil
}
