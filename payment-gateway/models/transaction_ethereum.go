package models

import (
	"github.com/jmoiron/sqlx"
)

type TxEthereum struct {
	ID       int     `json:"id"`
	Sender   string  `json:"sender"`
	Txhash   string  `json:"txhash"`
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

// NewTxEthereum ..
func NewTxEthereum(sender, txhash string, value float64) *TxEthereum {
	return &TxEthereum{
		Sender: sender,
		Txhash: txhash,
		Value:  value,
	}
}

// Add create user
func (c *TxEthereum) Add(db *sqlx.DB) (int, error) {
	var id int
	value := c.Value
	err := db.QueryRow(`INSERT INTO transaction_ethereum (sender, txhash, amount) VALUES($1, $2, $3) RETURNING id`, c.Sender, c.Txhash, value).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
