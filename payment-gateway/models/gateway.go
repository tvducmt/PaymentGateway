package models

import (
	"github.com/jmoiron/sqlx"
)

// FindGatewayByBankName Tìm kiếm id của gateway khi
func FindGatewayByBankName(db *sqlx.DB, bankName string) (int, error) {
	var bankID int
	err := db.Get(&bankID, "SELECT id FROM payment_gateway WHERE name = $1", bankName)
	return bankID, err
}
