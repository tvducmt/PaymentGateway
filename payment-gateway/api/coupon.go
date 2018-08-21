package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
	"gitlab.com/rockship/payment-gateway/models"
)

// CouponWithNewTXProcess give user_id and value
// check available balance =>user
// create transaction
// create receipt
// create coupon
// CouponWithNewTXProcess
/*
func (e *Env) CouponWithNewTXProcess(w http.ResponseWriter, r *http.Request) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Info("POST /coupon/new-transaction > CouponWithNewTXProcess Error")
		return StatusError{500, err, nil}
	}
	var requestData struct {
		UserID int     `json:"user_id"`
		Value  float64 `json:"value"`
	}
	err = json.Unmarshal(b, &requestData)
	if err != nil {
		logrus.Info("POST /coupon/new-transaction > Unmarshal Error")
		return StatusError{500, err, nil}
	}
	var responseData struct {
		Error  string
		Coupon *models.Coupon
	}

	user := models.Account{}
	userInfo, err := user.FindUserInfoByID(e.DB, requestData.UserID)
	if err != nil {
		logrus.Info("POST /coupon/new-transaction > FindUserInfoByID Error")
		return StatusError{201, err, []byte(`{"error": "User not exist"}`)}
	}
	if userInfo.Balance < requestData.Value {
		responseData.Error = "Not enough money to make a transaction"
		output, err := json.Marshal(responseData)
		if err != nil {
			logrus.Info("POST /coupon/new-transaction > Marshal Error")
			return StatusError{500, err, nil}
		}
		return StatusError{200, nil, output}
	}
	txHandle, err := e.DB.Begin()
	if err != nil {
		logrus.Info("POST /coupon/new-transaction > e.DB.Begin Error")
		return StatusError{500, err, nil}
	}
	tx := models.Transaction{}
	txID, err := tx.AddWithCoupon(txHandle, requestData.UserID)
	if err != nil {
		responseData.Error = err.Error()
		output, err := json.Marshal(responseData)
		if err != nil {
			logrus.Info("POST /coupon/new-transaction > Marshal Error")
			return StatusError{500, err, nil}
		}
		return StatusError{200, nil, output}
	}
	rc := models.Receipt{}
	receiptID, err := rc.AddWithCoupon(txHandle, requestData.Value)
	if err != nil {
		responseData.Error = err.Error()
		output, err := json.Marshal(responseData)
		if err != nil {
			logrus.Info("POST /coupon/new-transaction > Marshal Error")
			return StatusError{500, err, nil}
		}
		return StatusError{200, nil, output}
	}
	err = tx.UpdateTXWithCoupon(txHandle, txID, receiptID)
	if err != nil {
		responseData.Error = err.Error()
		output, err := json.Marshal(responseData)
		if err != nil {
			logrus.Info("POST /coupon/new-transaction > Marshal Error")
			return StatusError{500, err, nil}
		}
		return StatusError{200, nil, output}
	}
	cp := models.Coupon{}
	coupon, err := cp.CreateCouponWithNewReceipt(txHandle, requestData.UserID, txID, requestData.Value)
	if err != nil {
		responseData.Error = err.Error()
		output, err := json.Marshal(responseData)
		if err != nil {
			logrus.Info("POST /coupon/new-transaction > Marshal Error")
			return StatusError{500, err, nil}
		}
		return StatusError{200, nil, output}
	}
	responseData.Error = ""
	responseData.Coupon = coupon
	output, err := json.Marshal(responseData)
	if err != nil {
		logrus.Info("POST /coupon/new-transaction > Marshal Error")
		return StatusError{500, err, nil}
	}
	txHandle.Commit()
	return StatusError{200, nil, output}

}
*/

// GetCouponInfo swagger:route GET /user/coupons Coupon GetCouponInfo
//
// 	Lấy danh sách tất cả các coupon mà user đang sở hữu
//
// 	Responses:
//		200: DocResCouponInfo
//		500: DocResDefault
func (e *Env) GetCouponInfo(w http.ResponseWriter, r *http.Request) error {
	userOauth, err := GetUserByOauthContext(r)
	if err != nil {
		logrus.Info("GET /user/coupons > GetUserByOauthContext Error: ", err)
		return StatusError{201, err, []byte(`{"error": "Unauthorized"}`)}
	}
	userID := userOauth.ID

	var responseData struct {
		Error   string          `json:"error"`
		Coupons []models.Coupon `json:"coupons"`
	}
	tx := models.Transaction{}

	listTxID, err := tx.FindTxIDByUserID(e.DB, userID)
	if err != nil {
		logrus.Info("GET /user/coupons > FindTxIDByUserID Error: ", err)
		return StatusError{500, err, nil}
	}
	cp := models.Coupon{}
	listCoupon, err := cp.FindCouponsByTxID(e.DB, listTxID)

	if err != nil {
		log.Println(err)
		responseData.Error = err.Error()
		output, err := json.Marshal(responseData)
		if err != nil {
			logrus.Info("GET /user/coupons > Marshal Error: ", err)
			return StatusError{500, err, nil}
		}
		return StatusError{200, nil, output}
	}
	// sort list coupon by value
	// InsertionSort(listCoupon)
	responseData.Error = ""
	responseData.Coupons = listCoupon
	output, err := json.Marshal(responseData)
	if err != nil {
		logrus.Info("GET /user/coupons > Marshal Error: ", err)
		return StatusError{500, err, nil}
	}
	return StatusError{200, nil, output}
}
