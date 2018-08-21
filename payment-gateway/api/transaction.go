package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/stripe/stripe-go"

	"github.com/stripe/stripe-go/charge"
	"gitlab.com/rockship/payment-gateway/models"
)

type reqReceiveReceipt struct {
	Timestamp   string `json:"timestamp"`
	RawReceipt  string `json:"receipt"`
	Phonenumber string `json:"phonenumber"`
}

type responseData struct {
	Error  string
	TxCode string
}

// BankTransferProcess swagger:route POST /bank-transfer Transaction BankTransferProcess
//
// 	Thực hiện quá trình giao dịch bằng hình thức Bank Transfer cho user
//
// 	Responses:
//		500: DocResDefault
func (e *Env) BankTransferProcess(w http.ResponseWriter, r *http.Request) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return StatusError{500, err, nil}
	}
	var data struct {
		UserID    int `json:"user_id"`
		MethodID  int `json:"method_id"`
		GatewayID int `json:"gateway_id"`
	}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return StatusError{500, err, nil}
	}

	userOauth, err := GetUserByOauthContext(r)
	if err != nil || userOauth.ID != data.UserID {
		return StatusError{500, errors.New("Unauthorized"), nil}
	}

	tx := models.NewTransaction(data.UserID, data.MethodID, data.GatewayID)
	txcode, err := tx.Add(e.DB)
	if err != nil {
		return StatusError{500, nil, nil}
	}

	res := responseData{
		Error:  "",
		TxCode: txcode,
	}
	output, err := json.Marshal(res)
	if err != nil {
		return StatusError{500, err, nil}
	}
	return StatusError{200, nil, output}
}

// EthereumTransferProcess..
func (e *Env) EthereumTransferProcess(w http.ResponseWriter, r *http.Request) error {
	var data struct {
		UserID    int `json:"user_id"`
		MethodID  int `json:"method_id"`
		GatewayID int `json:"gateway_id"`
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return StatusError{500, err, nil}
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return StatusError{500, err, nil}
	}

	tx := models.NewTransaction(data.UserID, data.MethodID, data.GatewayID)
	txcode, err := tx.Add(e.DB)
	if err != nil {
		return StatusError{500, err, nil}
	}
	res := responseData{
		Error:  "",
		TxCode: txcode,
	}
	output, err := json.Marshal(res)
	if err != nil {
		return StatusError{500, err, nil}
	}
	return StatusError{200, nil, output}
}

// UpdateTransactionTimeOut swagger:route  PUT /transactions/timeout Transaction UpdateTransactionTimeOut
//
// 	Giám sát và cập nhật lại trạng thái transaction khi quá thời hạn để được xử lí
//
// 	Responses:
//		200: DocResUpdateTxTimeOut
//		500: DocResDefault
func (e *Env) UpdateTransactionTimeOut(w http.ResponseWriter, r *http.Request) error {
	tx := models.Transaction{}
	err := tx.CheckTransactionTimeOut(e.DB)
	if err != nil {
		return StatusError{500, err, nil}
	}
	result := []byte("Update successfull")
	return StatusError{200, nil, result}
}

// HistoryTransaction swagger:route  POST /transactions Transaction HistoryTransaction
//
// 	Trả về thông tin lịch sử tất cả các giao dịch hiện có của user
//
// 	Responses:
//		500: DocResDefault
func (e *Env) HistoryTransaction(w http.ResponseWriter, r *http.Request) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return StatusError{500, err, nil}
	}

	var data struct {
		UserID int `json:"userid"`
	}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return StatusError{500, err, nil}
	}
	txs, err := models.FindTransactionByUserID(e.DB, data.UserID)
	result := struct {
		Data  []models.TransactionUser `json:"data"`
		Error bool
	}{
		txs,
		false,
	}
	if err != nil {
		fmt.Println("Err FindTransactionByUserID POST /dashboard: ", err)
		result.Error = true
		result.Data = nil
	}

	jsnResult, err := json.Marshal(result)
	if err != nil {
		log.Println("POST Dashboard: ", err)
		return StatusError{500, err, nil}
	}
	return StatusError{200, nil, jsnResult}
}

// GetTransactionDetail swagger:route  GET /transaction Transaction GetTransactionDetail
//
// 	Lấy thông tin chi tiết của một giao dịch dựa vào id do user cung cấp
//
// 	Responses:
//		200: DocResTxSuccess
//		201: DocResTxMsg
//		500: DocResDefault
func (e *Env) GetTransactionDetail(w http.ResponseWriter, r *http.Request) error {
	_, err := GetUserByOauthContext(r)
	if err != nil {
		jsnResult := ([]byte(`{"error": "Unauthorized"}`))
		return StatusError{201, err, jsnResult}
	}

	r.ParseForm()
	txStrID := r.FormValue("tx_id")
	txID, err := strconv.Atoi(txStrID)
	if err != nil {
		log.Println("GetTransactionDetail > Atoi: ", err)
		jsnResult := ([]byte(`{"error": "Transaction ID Invalid"}`))
		return StatusError{201, err, jsnResult}
	}

	tx, err := models.FindFullTransactionByID(e.DB, txID)
	if err != nil {
		log.Println("GetTransactionDetail > FindFullTransactionByID: ", err)
		jsnResult := ([]byte(`{"error": "Transaction Not Exist!"}`))
		return StatusError{201, err, jsnResult}
	}

	data := struct {
		Transaction models.TransactionUser `json:"transaction"`
		MetaCredit  *stripe.Charge         `json:"charge"`
	}{
		tx[0],
		nil,
	}
	if tx[0].ChargeID.Valid != false {
		c, err := charge.Get(tx[0].ChargeID.String, nil)
		if err != nil {
			jsnResult := ([]byte(`{"error": "Someting went wrong. Please Try again!"}`))
			return StatusError{201, err, jsnResult}
		}
		data.MetaCredit = c
	}

	b, err := json.Marshal(data)
	if err != nil {
		log.Println("GetTransactionDetail > Marshal: ", err)
		jsnResult := ([]byte(`{"error": "Someting went wrong. Please Try again!"}`))
		return StatusError{201, err, jsnResult}
	}
	return StatusError{200, nil, b}
}
