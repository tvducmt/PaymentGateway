package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
	"gitlab.com/rockship/payment-gateway/models"
)

// StandardResponse ...
type StandardResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ErrCookieUserIDNotExist ...
var ErrCookieUserIDNotExist = errors.New("Cookie UserID Not Exist")

// ErrCookitExistButUserNotExist ...
var ErrCookitExistButUserNotExist = errors.New("Cookie Exist But User Not Exist")

func getCustomerByToken(token string) (*stripe.Customer, error) {

	params := &stripe.CustomerParams{
		Description: stripe.String("Customer for elizabeth.williams@example.com"),
	}
	err := params.SetSource(token) // obtained with Stripe.js
	if err != nil {
		log.Println("SetSource err: ", err)
		return nil, err
	}

	cus, err := customer.New(params)
	if err != nil {
		log.Println("Customer err: ", err)
		return nil, err
	}
	return cus, nil
}

func chargeByCustomerID(name, cusID string, amount int64) (*stripe.Charge, error) {
	paramsCard := &stripe.CardParams{
		Customer: stripe.String(cusID),
	}
	// cardUser, err := card.New(paramsCard)
	sp := &stripe.SourceParams{
		Card: paramsCard,
	}
	chParams := &stripe.ChargeParams{
		Amount:      stripe.Int64(amount),
		Currency:    stripe.String(string(stripe.CurrencyVND)),
		Description: stripe.String("Customer name: " + name),
		Source:      sp,
		Customer:    stripe.String(cusID),
	}
	chParams.SetSource(paramsCard)
	ch, err := charge.New(chParams)
	if err != nil {
		return nil, nil
	}
	// fmt.Println(ch.ID)
	return ch, nil
}

// InsertionSort ...
func InsertionSort(array []models.Coupon) {
	for i := 1; i < len(array); i++ {
		for j := i; j > 0 && array[j].Value > array[j-1].Value; j-- {
			array[j], array[j-1] = array[j-1], array[j]
		}
	}
}

// CreateTxCreditCard ...
func CreateTxCreditCard(db *sqlx.DB, userID, methodID int, chargeID string) error {
	var data struct {
		UserID   int    `json:"user_id"`
		MethodID int    `json:"method_id"`
		ChargeID string `json:"charge_id"`
	}
	data.UserID = userID
	data.MethodID = methodID
	data.ChargeID = chargeID

	tx := models.NewCreditCard(data.UserID, data.MethodID, data.ChargeID)

	// Tạo transaction với charge id tương ứng khi charge thành công
	err := tx.AddCreditCard(db)
	if err != nil {
		log.Println("CreateTxCreditCard > AddCreditCard: ", err)
		return err
	}

	c, err := charge.Get(chargeID, nil)
	if err != nil {
		log.Println("Helper > Get Charge error: ", err)
		return err
	}
	// c.Amount
	transactionSQL, err := db.Begin()
	if err != nil {
		fmt.Println("Func CreateTxCreditCard > AddCreditCard : ", err)
		return err
	}

	// Tạo coupon với giá trị tương đương
	cp, err := models.NewCoupon(transactionSQL, tx.ID, float64(c.Amount), string(c.Currency))
	if err != nil {
		transactionSQL.Rollback()
		log.Println("Error create coupon ", err)
		return err
	}
	coupon, err := cp.Add(transactionSQL)
	if err != nil {
		transactionSQL.Rollback()
		log.Println("Func CreateTxCreditCard > Add Coupon Err: ", err)
		return err
	}
	_ = coupon

	err = transactionSQL.Commit()
	if err != nil {
		transactionSQL.Rollback()
		log.Println("Func CreateTxCreditCard Commit Err: ", err)
		return err
	}
	err = models.UpdateSuccess(db, tx.ID)
	if err != nil {
		log.Println("Helper > UpdateSuccess Tx Credit Card: ", err)
		return err
	}
	return nil
}

// GetUserByOauthContext ...
func GetUserByOauthContext(r *http.Request) (*OauthUser, error) {
	userOauth, ok := r.Context().Value("user").(OauthUser)
	if !ok {
		return nil, errors.New("Unauthorized")
	}

	return &userOauth, nil
}

func CovertStringToFloat(amount string) float64 {
	amountString := strings.Replace(amount, ",", "", -1)
	amountMoney, err := strconv.Atoi(amountString)
	if err != nil {
		return 0
	}
	amountFloat := float64(amountMoney)
	return amountFloat
}

func CallApiVerifyAccessToken(accessToken string) (*http.Response, *OauthUser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", os.Getenv("OAUTH_VERIFY"), nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var user OauthUser
	err = json.Unmarshal(b, &user)
	if err != nil {
		return nil, nil, err
	}
	return resp, &user, nil
}

// CallApiOauthToken ...
func CallApiOauthToken(email, password string) (*http.Response, *ResOauthService, error) {
	params := url.Values{}
	params.Set("grant_type", "password")
	params.Set("username", email)
	params.Set("password", password)
	postData := strings.NewReader(params.Encode())

	client := &http.Client{}

	req, err := http.NewRequest("POST", os.Getenv("OAUTH_URI"), postData)

	if err != nil {
		return nil, nil, err
	}
	header := os.Getenv("CLIENT_ID") + ":" + os.Getenv("CLIENT_SECRET")
	base64Header := base64.StdEncoding.EncodeToString([]byte(header))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64Header)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		return nil, nil, err
	}

	var respOauth ResOauthService
	err = json.Unmarshal(b, &respOauth)
	if err != nil {
		return nil, nil, err
	}

	return resp, &respOauth, nil
}
