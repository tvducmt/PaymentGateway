package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gitlab.com/rockship/payment-gateway/models"
)

// ReceiveReceipt swagger:route  POST /secret/receipt Receipt ReceiveReceipt
//
// 	Tiếp nhận Receipt từ thiết bị di động, xử lí và trích xuất những thông tin cần thiết dựa vào Transaction Code có trong Receipt
//	. Nếu Receipt gửi lên hợp lệ, tạo một Coupon có giá trị tương tương cho user tương ứng.
//
// 	Responses:
//		200: DocReceiveReceiptSuccess
//		500: DocResDefault
func (e *Env) ReceiveReceipt(w http.ResponseWriter, r *http.Request) error {
	receiptJsn := reqReceiveReceipt{}

	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return StatusError{500, err, nil}
	}
	var responseData struct {
		Error  string         `json:"error"`
		Coupon *models.Coupon `json:"coupon"`
	}
	err = json.Unmarshal(jsn, &receiptJsn)
	if err != nil {
		resByte, err := json.Marshal(StandardResponse{false, "Malicious Receipt."})
		if err != nil {
			return StatusError{500, err, nil}
		}
		// w.Write(resByte)
		// w.WriteHeader(http.StatusBadRequest)
		return StatusError{200, nil, resByte}
	}
	createAt, err := time.Parse("2006-01-02T15:04:05-07:00", receiptJsn.Timestamp)
	if err != nil {
		resByte, err := json.Marshal(StandardResponse{false, "Malicious Receipt. Wrong time format."})
		if err != nil {
			return StatusError{500, err, nil}
		}

		// w.Write(resByte)
		// w.WriteHeader(http.StatusBadRequest)
		return StatusError{200, nil, resByte}
	}
	txHandle, err := e.DB.Begin()
	if err != nil {
		return StatusError{500, err, nil}
	}
	receipt := models.NewReceipt(receiptJsn.RawReceipt, receiptJsn.Phonenumber, createAt)
	receiptID, err := receipt.Add(txHandle)
	if err != nil {
		txHandle.Rollback()
		resByte, err := json.Marshal(StandardResponse{false, "System error. Try again later."})
		if err != nil {
			fmt.Println("Marshal StandardResponse Error: ", err)
			return StatusError{500, err, nil}
		}
		return StatusError{200, nil, resByte}
	}
	rules, err := models.LoadRegex(e.DB)
	if err != nil {
		fmt.Println("LoadRegex Error: ", err)
		return StatusError{500, err, nil}
	}
	checkRule, err := models.CheckRule(receipt.RawReceipt, rules)
	if err != nil {
		fmt.Println("CheckRule Error: ", err)
		return StatusError{500, err, nil}
	}
	var txID int
	if checkRule != nil {
		txcode := checkRule["transaction_code"]
		tx := models.Transaction{}
		txID, err = tx.UpdateTXByCode(txHandle, receiptID, txcode)
		if err != nil {
			txHandle.Rollback()
			fmt.Println("ReceiveReceipt > UpdateTXByCode Error: ", err)
			return StatusError{500, err, nil}
		}
		err = models.UpdateReceiptByreceiptID(txHandle, receiptID, checkRule)
		if err != nil {
			txHandle.Rollback()
			fmt.Println("ReceiveReceipt > UpdateReceiptByreceiptID Error: ", err)
			return StatusError{500, err, nil}
		}

		// input iduser, idtransaction
		// Create coupon
		amount := CovertStringToFloat(checkRule["amount"])

		coupon, err := models.NewCoupon(txHandle, txID, amount, "vnd")
		if err != nil {
			log.Println("Create coupon failed", err)
			return StatusError{500, err, nil}
		}
		coupon, err = coupon.Add(txHandle)
		if err != nil {
			txHandle.Rollback()
			log.Println("Add failed", err)
			return StatusError{500, err, nil}
		}
		responseData.Error = ""
		responseData.Coupon = coupon
		output, err := json.Marshal(responseData)
		if err != nil {
			return StatusError{500, err, nil}
		}

		txHandle.Commit()
		return StatusError{200, nil, output}
	}

	return StatusError{500, err, nil}
}
