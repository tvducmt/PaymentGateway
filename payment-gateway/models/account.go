package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
)

// swagger:model Account
type Account struct {
	ID           int            `json:"id" db:"id"`
	Email        string         `json:"email" db:"email"`
	Passphrase   string         `json:"passphrase" db:"passphrase"`
	Fullname     sql.NullString `json:"fullname" db:"fullname"`
	Phone        sql.NullString `json:"phone" db:"phone"`
	OtpEnable    bool           `json:"otpenable" db:"otpenable"`
	OtpSecretKey sql.NullString `json:"otpsecretkey" db:"otpsecretkey"`
}

type FogotPwRequest struct {
	ID       int       `json:"id"`
	UID      int       `json:"uid"`
	Token    string    `json:"token"`
	CreateAt time.Time `json:"time"`
}

// swagger:model  UserInfo
type UserInfo struct {
	// ID định danh cho user
	//
	// min: 1
	UserID int `json:"userid"`

	// Địa chỉ Email của user có định danh id
	Email string `json:"email"`

	// Số điện thoại của user tương ứng
	Phone string `json:"phone"`

	// Họ và tên do user cung cấp
	Fullname string `json:"fullname"`

	OtpEnable bool `json:"otpenable"`

	OtpSecretKey string `json:"otpsecretkey"`
}

// swagger:model  Password
type Password struct {
	// required: true
	// min: 1
	UserID int `json:"userid"`
	// required: true
	// min length: 6
	OldPassword string `json:"oldpassword"`
	// required: true
	// min length: 6
	NewPassword string `json:"newpassword"`
}

// NewAccount ...
func NewAccount(email, password string) *Account {
	var user Account
	user.Email = email
	user.Passphrase = password
	return &user
}

// NewReqForgotPw ...
func NewReqForgotPw(uid int, token string) *FogotPwRequest {
	var req FogotPwRequest
	req.UID = uid
	req.Token = token
	return &req
}

// UpdateUser ...
func (a *Account) UpdateUser(db *sqlx.DB, user *UserInfo) error {
	if user.OtpEnable && !a.OtpEnable {
		query := "UPDATE account SET fullname = $1, phone = $2, otpenable = $3, otpsecretkey= $4 WHERE id = $5"
		_, err := db.Exec(query, user.Fullname, user.Phone, user.OtpEnable, user.OtpSecretKey, user.UserID)
		if err != nil {
			return err
		}
	} else {
		query := "UPDATE account SET fullname = $1, phone = $2, otpenable = $3 WHERE id = $4"
		_, err := db.Exec(query, user.Fullname, user.Phone, user.OtpEnable, user.UserID)
		if err != nil {
			return err
		}
	}

	return nil
}

//UpdatePassword ...
func (a *Account) UpdatePassword(db *sqlx.DB, pw *Password) error {
	user, err := FindUserByID(db, pw.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("User no exist")
	}
	_, err = FindUserByEmailPassword(db, user.Email, pw.OldPassword)
	if err != nil {
		return errors.New("Old password not incorrect, Please try again")
	}

	query := "UPDATE account SET  passphrase = crypt($1, passphrase) WHERE id = $2"
	_, err = db.Exec(query, pw.NewPassword, pw.UserID)
	if err != nil {
		return errors.New("Update account error")
	}
	return nil
}

//ResetPassword ...
func (a *Account) ResetPassword(db *sql.Tx, newPassword string) error {
	query := "UPDATE account SET passphrase = crypt($1, passphrase) WHERE id = $2"
	_, err := db.Exec(query, newPassword, a.ID)
	if err != nil {
		return errors.New("Update account error")
	}
	return nil
}

// FindUserInfoByID find user by id
func (a *Account) FindUserInfoByID(db *sqlx.DB, id int) (*UserInfo, error) {
	var acc Account
	err := db.QueryRow(`SELECT id, email, fullname, phone, otpenable FROM account WHERE id = $1 `, id).Scan(&acc.ID, &acc.Email, &acc.Fullname, &acc.Phone, &acc.OtpEnable)
	if err != nil {
		return nil, err
	}
	return &UserInfo{
		UserID:    acc.ID,
		Email:     acc.Email,
		Fullname:  acc.Fullname.String,
		Phone:     acc.Phone.String,
		OtpEnable: acc.OtpEnable,
	}, nil
}

// FindUserByID ...
func FindUserByID(db *sqlx.DB, id int) (*Account, error) {
	var account Account
	err := db.Get(&account, `SELECT * FROM account WHERE id = $1`, id)
	if err != nil {
		fmt.Println("FindUserByID: ", err)
		return nil, err
	}
	account.Passphrase = "" // For security
	return &account, nil
}

// FindUserByEmail ...
func FindUserByEmail(db *sqlx.DB, email string) (*Account, error) {
	var account Account

	err := db.Get(&account, `SELECT * FROM account WHERE email = $1`, email)
	if err != nil {
		return nil, err
	}
	account.Passphrase = "" // For security
	return &account, nil
}

// FindUserByEmailPassword ...
func FindUserByEmailPassword(db *sqlx.DB, email, password string) (*Account, error) {
	var account Account
	err := db.Get(&account, `SELECT * FROM account WHERE email = $1 AND passphrase = crypt($2, passphrase)`, email, password)
	if err != nil {
		return nil, err
	}
	account.Passphrase = "" // For security
	return &account, nil
}

// Add create user
func (a *Account) Add(db *sqlx.DB) (int, error) {
	var id int
	err := db.QueryRow(`INSERT INTO account (email, passphrase) VALUES($1, crypt($2, gen_salt('bf', 8))) RETURNING id`, a.Email, a.Passphrase).Scan(&id)
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

// Add request forgot password
func (fw *FogotPwRequest) Add(tx *sql.Tx) error {
	err := tx.QueryRow(`INSERT INTO password_change_requests (user_id, token) VALUES($1, crypt($2, gen_salt('bf', 8))) RETURNING id`, fw.UID, fw.Token).Scan(&fw.ID)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			if pqErr.Code.Name() == "unique_violation" {
				fw.ID = -1
				return nil
			}
		}
		return err
	}
	return nil
}

// FindReqForgotPwByToken ...
func FindReqForgotPwByToken(db *sql.Tx, token string) (*FogotPwRequest, error) {
	var req FogotPwRequest
	err := db.QueryRow(`SELECT id, user_id, token, create_at 
		FROM password_change_requests 
		WHERE token = crypt($1, token) AND DATE_PART('Day',now() - create_at) < 1`, token).Scan(&req.ID, &req.UID, &req.Token, &req.CreateAt)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`DELETE FROM password_change_requests WHERE token = crypt($1, token)`, token)
	if err != nil {
		return nil, err
	}

	req.Token = ""
	return &req, nil
}
