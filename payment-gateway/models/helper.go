package models

import (
	"crypto/rand"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	CodeQuantity = 450 //450000000
)

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(n int) (string, error) {
	const letters = "ABCDFGHIJKLMNPQRSTVWXYZ0123456789"
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

//InsertCodeTable ..
func InsertCodeTable(db *sqlx.DB, code string) error {
	query := `INSERT INTO codetable (code) VALUES ($1)`
	_, err := db.Exec(query, code)
	if err != nil {
		return err
	}
	return nil
}

// GenarateCodes ...
func GenarateCodes(db *sqlx.DB) error {
	for i := 0; i < CodeQuantity; i++ {
		code, err := GenerateRandomString(6)
		if err != nil {
			fmt.Println("Error generate string", err)
			return err
		}
		err = InsertCodeTable(db, code)
		if err != nil {
			pqErr, ok := err.(*pq.Error)
			if ok {
				if pqErr.Code.Name() == "unique_violation" {
					continue
				}
			}
			fmt.Println("Error insert code", err)
			return err
		}
	}
	for i := 0; i < CodeQuantity; i++ {
		code, err := GenerateRandomString(12)
		if err != nil {
			fmt.Println("Error generate string", err)
			return err
		}
		err = InsertCodeTable(db, code)
		if err != nil {
			pqErr, ok := err.(*pq.Error)
			if ok {
				if pqErr.Code.Name() == "unique_violation" {
					continue
				}
			}
			fmt.Println("Error insert code", err)
			return err
		}
	}
	return nil
}
