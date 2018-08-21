package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"gitlab.com/rockship/payment-gateway/models"
)

// GetCardsStripe swagger:route GET /customer/cards Customer GetCardsStripe
//
// 	Lấy danh sách tất cả các Credit Card mà user đã lưu ở các giao dịch trước đó
//
// 	Responses:
//		200: DecResCustomer
//		500: DocResDefault
func (e *Env) GetAllCustomerStripe(w http.ResponseWriter, r *http.Request) error {
	stripe.Key = os.Getenv("STRIPE_SK")
	userOauth, err := GetUserByOauthContext(r)
	if err != nil {
		logrus.Info("GET /customer/cards > GetUserByOauthContext Error")
		return StatusError{500, err, nil}
	}
	userID := userOauth.ID
	if err != nil {
		logrus.Info("GET /customer/cards > Int to Float64 Error")
		return StatusError{500, err, nil}
	}
	cus, err := models.FindCustomerByUserID(e.DB, userID)
	if err != nil {
		logrus.Info("GET /customer/cards > FindCustomerByUserID Error")
		return StatusError{500, err, nil}
	}
	cusUsers := make([]models.CustomerUser, len(cus))

	for k, v := range cus {
		cusUsers[k].ID = v.ID
		cusUsers[k].UserID = v.UserID
		// cusUsers[k].CusID = v.CusStripeID

		c, err := customer.Get(v.CusStripeID, nil)
		if err != nil {
			logrus.Info("GET /customer/cards > Get Charge Error")
			return StatusError{500, err, nil}
		}
		cusUsers[k].Last4 = (*(*c).Sources).Data[0].Card.Last4
		cusUsers[k].ExpMonth = (*(*c).Sources).Data[0].Card.ExpMonth
		cusUsers[k].ExpYear = (*(*c).Sources).Data[0].Card.ExpYear
		cusUsers[k].AddressZip = (*(*c).Sources).Data[0].Card.AddressZip
		cusUsers[k].Brand = (*(*c).Sources).Data[0].Card.Brand
		cusUsers[k].Funding = (*(*c).Sources).Data[0].Card.Funding
		cusUsers[k].Fingerprint = (*(*c).Sources).Data[0].Card.Fingerprint
	}

	b, err := json.Marshal(cusUsers)
	if err != nil {
		logrus.Info("GET /customer/cards > Marshal Error")
		return StatusError{500, err, nil}
	}
	return StatusError{200, nil, b}
}
