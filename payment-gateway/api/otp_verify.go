package api

import (
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	otp "github.com/hgfischer/go-otp"
	"github.com/sirupsen/logrus"
	"gitlab.com/rockship/payment-gateway/helper"
	"gitlab.com/rockship/payment-gateway/models"
)

//VerifyOtp ...
func (e *Env) VerifyOtp(w http.ResponseWriter, r *http.Request) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Info("POST /register > RegisterProcess > ReadAll Error: ", err)
		return StatusError{500, err, nil}
	}
	var req struct {
		UserID             int    `json:"user_id"`
		OtpPin             string `json:"otp_pin"`
		AccesstokenEncrypt string `json:"accesstoken_encrypted"`
	}
	err = json.Unmarshal(b, &req)
	if err != nil {
		logrus.Info("POST /register > RegisterProcess > ReadAll Error: ", err)
		return StatusError{500, err, nil}
	}

	acc, err := models.FindUserByID(e.DB, req.UserID)
	if err != nil {
		logrus.Info("POST /user/otp > FindUserByID: ", err)
		return StatusError{500, err, nil}
	}
	if acc.OtpEnable {
		totp := &otp.TOTP{Secret: acc.OtpSecretKey.String, IsBase32Secret: true}
		if totp.Verify(req.OtpPin) {
			key := []byte(os.Getenv("CLIENT_SECRET"))
			key = key[:32]

			ciphertextDecode, err := base32.StdEncoding.DecodeString(req.AccesstokenEncrypt)
			if err != nil {
				return StatusError{500, err, nil}
			}
			AccesstokenDecrypt, err := helper.AESDecrypt(ciphertextDecode, key)
			if err != nil {
				logrus.Println("Error Decrypt", err)
			}
			respVerify, oauthUser, err := CallApiVerifyAccessToken(string(AccesstokenDecrypt))
			if err != nil || respVerify.StatusCode != 200 {
				fmt.Println("CallApiVerifyAccessToken ", err)
				return StatusError{500, err, nil}
			}
			data := struct {
				Status int       `json:"status"`
				Result OauthUser `json:"result"`
				Token  string    `json:"token"`
			}{
				200,
				*oauthUser,
				string(AccesstokenDecrypt),
			}

			b, err := json.Marshal(data)
			if err != nil {
				return StatusError{500, err, nil}
			}
			return StatusError{200, nil, b}
		}
	}
	return StatusError{200, nil, []byte(`{"error": "OTP Code Invalid"}`)}
}
