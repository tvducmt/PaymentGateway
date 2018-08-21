package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"gitlab.com/rockship/payment-gateway/models"
)

// CreditCardProcess swagger:route POST /credit-card Transaction CreditCardProcess
//
// 	CreditCardProcess
//
// 	Thực hiện quá trình giao dịch bằng hình thức Credit Card cho user. Charge thành công tạo ngay một Coupon có giá trị tương đương cho người dùng.
//
// 	Responses:
//		200:
//		500: DocResDefault
func (e *Env) CreditCardProcess(w http.ResponseWriter, r *http.Request) error {
	stripe.Key = os.Getenv("STRIPE_SK")
	userOauth, err := GetUserByOauthContext(r)
	if err != nil {
		return StatusError{500, err, nil}
	}
	userID := userOauth.ID

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("CreditCardProcess: ", err)
		return StatusError{500, err, nil}
	}

	var data struct {
		Name     string `json:"name"`
		Amount   int64  `json:"amount"`
		Token    string `json:"token"`
		ModeCard string `json:"select_mode_card"`
		CusID    int    `json:"select_cus_id"`
		Remember bool   `json:"remember_me"`
	}
	err = json.Unmarshal(b, &data)
	if err != nil {
		log.Println("Err Unmarshal CreditCard: ", err)
		return StatusError{500, err, nil}
	}
	var cusStrID string
	if data.ModeCard == "new" {
		// TODO: Get customer by token
		cus, err := getCustomerByToken(data.Token)
		if err != nil {
			log.Println("Err getCustomerByToken: ", err)
			return StatusError{500, err, nil}
		}

		if data.Remember {
			cusDB := models.NewCustomer(userID, cus.ID)
			c, err := customer.Get(cus.ID, nil)
			if err != nil {
				fmt.Println("CreditCardProcess > NewCustomer: ", err)
				return StatusError{500, err, nil}
			}
			cusDB.Fingerprint = (*(*c).Sources).Data[0].Card.Fingerprint

			cusDB.ID, err = cusDB.AddCustomerToUser(e.DB)
			// Skip err if cusDB.ID = -1 => Credit Card existing in DB. Continue Charge
			if err != nil {
				log.Println("Err AddCustomerToUser: ", err)
				return StatusError{500, err, nil}
			}
		}

		cusStrID = cus.ID
	} else if data.ModeCard == "old" {
		// TODO: Get customer by DB
		customer, err := models.FindCustomerByID(e.DB, data.CusID)
		if err != nil {
			log.Println("Err FindCustomerByID Credit: ", err)
			return StatusError{500, err, nil}
		}
		cusStrID = customer.CusStripeID
	}

	ch, err := chargeByCustomerID(data.Name, cusStrID, data.Amount)
	if err != nil {
		log.Println("Charge err: ", err)
		return StatusError{500, err, nil}
	}
	if ch != nil {
		err = CreateTxCreditCard(e.DB, userID, 2, ch.ID)
		if err != nil {
			fmt.Println("CreditCardProcess > CreateTxCreditCard: ", err)
			return StatusError{500, err, nil}
		}
		return StatusError{200, nil, []byte(`{"error": ""}`)}
	}

	return StatusError{500, err, nil}
}
