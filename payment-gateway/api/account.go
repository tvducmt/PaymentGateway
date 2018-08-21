package api

import (
	"crypto"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/sirupsen/logrus"

	"gitlab.com/rockship/payment-gateway/helper"
	"gitlab.com/rockship/payment-gateway/models"
	"gitlab.com/rockship/payment-gateway/services"
)

var (
	// TokenOauth ...
	TokenOauth   *ResOauthService
	errBadFormat = errors.New("Invalid format")
	emailRegexp  = regexp.MustCompile(`^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
	pwRegexp     = regexp.MustCompile(`.{6,32}`)
)

// swagger:model  AccountAuth
type AccountAuth struct {
	// required: true
	Email string `json:"email"`

	// required: true
	Password string `json:"password"`
}

type ResOauthService struct {
	AccessToken  string `json:"access_token"`
	ExpriesIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type OauthUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// Message ...
type message struct {
	Status    string    `json:"status"`
	Error     string    `json:"error"`
	Token     string    `json:"token"`
	Result    OauthUser `json:"result"`
	OtpEnable bool      `json:"otpenable"`
}

func validateReqData(email, pw string) error {
	if !emailRegexp.MatchString(email) || !pwRegexp.MatchString(pw) {
		return errBadFormat
	}
	return nil
}

// LoginProcess swagger:route POST /login Users Login
//
// 	Xử lý quá trình đăng nhập của user
//
// 	Responses:
//		200: DocResLogin
//		500: DocResDefault
func (e *Env) LoginProcess(w http.ResponseWriter, r *http.Request) error {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logrus.Info("POST /login > LoginProcess > ReadAll Error")
		return StatusError{500, err, nil}
	}

	var acc AccountAuth
	err = json.Unmarshal(b, &acc)
	if err != nil {
		logrus.Info("POST /login > LoginProcess > Unmarshal Error")
		return StatusError{500, err, nil}
	}

	var responseMsg message
	responseMsg.Status = "OK"

	resp, TokenOauth, err := CallApiOauthToken(acc.Email, acc.Password)

	if err != nil || resp.StatusCode != 200 {
		log.Println("NewRequest Oauth: ", err)
		return StatusError{500, errors.New("Unauthorized"), nil}
	}

	respVerify, user, err := CallApiVerifyAccessToken(TokenOauth.AccessToken)
	if err != nil || respVerify.StatusCode != 200 {
		responseMsg.Status = "ERROR"
		responseMsg.Error = "Email or Password wrong!"

		output, err := json.Marshal(responseMsg)
		if err != nil {
			logrus.Info("POST /login > LoginProcess > Marshal Error")
			return StatusError{500, err, nil}
		}
		return StatusError{200, nil, output}
	}

	account, err := models.FindUserByID(e.DB, user.ID)
	if err != nil {
		return StatusError{500, nil, []byte(`{"message": "User not exist"}`)}
	}

	if account.OtpEnable {
		key := []byte(os.Getenv("CLIENT_SECRET"))
		key = key[:32]
		ciphertext, err := helper.AESEncrypt([]byte(TokenOauth.AccessToken), key)
		if err != nil {
			fmt.Println("AESEncrypt: ", err)
			return StatusError{500, err, nil}
		}
		ciphertextEncode := base32.StdEncoding.EncodeToString(ciphertext)

		responseMsg.Result = *user
		responseMsg.Token = ciphertextEncode
		responseMsg.OtpEnable = true
		output, err := json.Marshal(responseMsg)

		if err != nil {
			return StatusError{500, err, nil}
		}
		return StatusError{200, nil, output}
	}
	// Check xem account nay co bat otp chua??
	// Neu bat thi phai lay trong
	responseMsg.Result = *user
	responseMsg.OtpEnable = false
	responseMsg.Token = TokenOauth.AccessToken
	output, err := json.Marshal(responseMsg)
	if err != nil {
		logrus.Info("POST /login LoginProcess Marshal Error")
		return StatusError{500, err, nil}
	}
	return StatusError{200, nil, output}
}

// RegisterProcess swagger:route POST /register Users Register
//
// 	Xử lý quá trình đăng ký tài khoản của user
//
// 	Responses:
//		200: DocResRegisterSuccess
//		201: DocResRegister
//		500: DocResDefault
func (e *Env) RegisterProcess(w http.ResponseWriter, r *http.Request) error {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logrus.Info("POST /register > RegisterProcess > ReadAll Error")
		return StatusError{500, err, nil}
	}
	var acc struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		PwConfirm string `json:"pwConfirm"`
	}
	err = json.Unmarshal(b, &acc)
	if err != nil {
		logrus.Info("POST /register > RegisterProcess > ReadAll Error")
		return StatusError{500, err, nil}
	}
	if acc.Password != acc.PwConfirm {
		return StatusError{501, errors.New("Password Not Match"), nil}
	}

	user := models.NewAccount(acc.Email, acc.Password)
	user.ID, err = user.Add(e.DB)
	if user.ID == -1 {
		return StatusError{501, errors.New("Existing Account"), nil}
	}
	if err != nil {
		log.Println("NewAccount Register Err: ", err)
		return StatusError{501, errors.New("Error establishing a database connection"), nil}
	}

	resp, TokenOauth, err := CallApiOauthToken(acc.Email, acc.Password)
	if err != nil || resp.StatusCode != 200 {
		logrus.Info("POST /register > RegisterProcess > Create AccessToken Error")
		return StatusError{500, err, nil}
	}

	respVerify, userOauth, err := CallApiVerifyAccessToken(TokenOauth.AccessToken)
	if err != nil || respVerify.StatusCode != 200 {
		fmt.Println("CallApiVerifyAccessToken ", err)
		return StatusError{500, err, nil}
	}

	dataRes := struct {
		Token  string    `json:"token"`
		Error  string    `json:"error"`
		Result OauthUser `json:"result"`
	}{
		TokenOauth.AccessToken,
		"",
		*userOauth,
	}
	output, err := json.Marshal(dataRes)
	if err != nil {
		logrus.Info("POST /register > RegisterProcess > Marshal Error")
		return StatusError{500, err, nil}
	}
	return StatusError{200, nil, output}
}

// UpdateUserInfo swagger:route PUT /user/info Users UpdateUserInfo
//
// 	Cập nhật thông tin cơ bản cho user. Chỉ cho phép thay đổi fullname và phone. Mọi thông tin khác gủi lên chỉ để xác minh user
//
// 	Responses:
//		200: DocResUpdateUserInfo
//		500: DocResDefault

func (e *Env) UpdateUserInfo(w http.ResponseWriter, r *http.Request) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Info("PUT /update-user > UpdateUserInfo > ReadAll Error")
		return StatusError{500, err, nil}
	}

	userInfo := models.UserInfo{}
	err = json.Unmarshal(b, &userInfo)
	if err != nil {
		logrus.Info("PUT /update-user > UpdateUserInfo > Unmarshal Error")
		return StatusError{500, err, nil}
	}

	userOauth, err := GetUserByOauthContext(r)
	if err != nil || userOauth.ID != userInfo.UserID {
		logrus.Info("PUT /update-user > GetUserByOauthContext Error")
		return StatusError{500, errors.New("Unauthorized"), nil}
	}

	acc, err := models.FindUserByID(e.DB, userOauth.ID)
	if err != nil {
		return StatusError{500, nil, []byte(`{"message": "User not exist"}`)}
	}
	if !acc.OtpEnable && userInfo.OtpEnable {
		var responseData struct {
			URLQrCode string `json:"urlQRCode"`
		}
		/*
			secretKey, err := models.GenerateRandomString(32)
			if err != nil {
				return StatusError{500, err, nil}
			}

			secretKey = base32.StdEncoding.EncodeToString([]byte(secretKey))
			if err != nil {
				return StatusError{500, err, nil}
			}
			userInfo.OtpSecretKey = secretKey
		*/
		keySize := crypto.SHA256.Size()
		key := make([]byte, keySize)
		_, err := rand.Read(key)
		if err != nil {
			return StatusError{500, err, nil}
		}
		secretKey := base32.StdEncoding.EncodeToString(key)
		userInfo.OtpSecretKey = secretKey
		err = acc.UpdateUser(e.DB, &userInfo)
		if err != nil {
			logrus.Info("PUT /update-user > UpdateUserInfo > UpdateUser Error")
			return StatusError{500, err, nil}
		}
		responseData.URLQrCode = fmt.Sprintf(`otpauth://totp/Rockship:%s?secret=%s`, acc.Email, secretKey)
		data, err := json.Marshal(responseData)
		if err != nil {
			return StatusError{500, err, nil}
		}
		return StatusError{200, nil, data}
	}

	err = acc.UpdateUser(e.DB, &userInfo)
	if err != nil {
		logrus.Info("PUT /update-user > UpdateUserInfo > UpdateUser Error")
		return StatusError{500, err, nil}
	}
	return StatusError{200, nil, []byte(`{"message": "Update successfull"}`)}
}

// UpdatePassword swagger:route PUT /user/password Users UpdatePassword
//
// 	Thay đổi mật khẩu cho user khi có yêu cầu
//
// 	Responses:
//		200: DocResUpdateUserInfo
//		500: DocResDefault
func (e *Env) UpdatePassword(w http.ResponseWriter, r *http.Request) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Info("PUT /user/password > UpdatePassword > ReadAll Error")
		return StatusError{500, err, nil}
	}
	pw := models.Password{}
	err = json.Unmarshal(b, &pw)
	if err != nil {
		logrus.Info("PUT /user/password > UpdatePassword > Unmarshal Error")
		return StatusError{500, err, nil}
	}

	userOauth, err := GetUserByOauthContext(r)
	if err != nil || userOauth.ID != pw.UserID {
		return StatusError{201, errors.New("Unauthorized"), []byte(`{"error": "Unauthorized"}`)}
	}

	acc := models.Account{}
	err = acc.UpdatePassword(e.DB, &pw)
	if err != nil {
		logrus.Info("PUT /user/password > UpdatePassword Error")
		return StatusError{501, err, nil}
	}
	return StatusError{200, nil, []byte(`{"message": "Change password success"}`)}
}

// GetUserInfo swagger:route GET /user/info Users GetUserInfo
//
// 	Lấy tất cả những thông tin cơ bản của user
//
// 	Responses:
//		200: DocResUserInfo
//		201: DocResUserMsg
//		500: DocResDefault
func (e *Env) GetUserInfo(w http.ResponseWriter, r *http.Request) error {
	userOauth, err := GetUserByOauthContext(r)
	if err != nil {
		logrus.Info("GET /userinfo > GetUserByOauthContext")
		return StatusError{201, err, []byte(`{"error": "Unauthorized"}`)}
	}

	userID := userOauth.ID

	user := models.Account{}
	userInfo, err := user.FindUserInfoByID(e.DB, userID)
	if err != nil {
		logrus.Info("GET /userinfo > FindUserInfoByID Error")
		return StatusError{201, err, []byte(`{"error": "User not exist"}`)}
	}

	resByteUser, err := json.Marshal(userInfo)
	if err != nil {
		logrus.Info("GET /userinfo > GetUserInfo Error")
		return StatusError{500, err, nil}
	}
	return StatusError{200, nil, resByteUser}
}

// ForgotPassword swagger:route POST /user/forgot-password Users ForgotPassword
//
// 	Lấy tất cả những thông tin cơ bản của user
//
// 	Responses:
//		500: DocResDefault
func (e *Env) ForgotPassword(w http.ResponseWriter, r *http.Request) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Info("POST /user/forgot-password > RegisterProcess > ReadAll Error")
		return StatusError{500, err, nil}
	}
	var dataReq struct {
		Email string `json:"email"`
	}
	err = json.Unmarshal(b, &dataReq)
	if err != nil {
		logrus.Info("POST /user/forgot-password > RegisterProcess > ReadAll Error")
		return StatusError{500, err, nil}
	}

	// FLow 1: Find user by email requested by client
	user, err := models.FindUserByEmail(e.DB, dataReq.Email)
	if err != nil {
		logrus.Info("POST /user/forgot-password > ForgotPassword > FindUserByEmail")
		return StatusError{201, err, []byte(`{"error": "Email not exist"}`)}
	}

	randomToken, err := models.GenerateRandomString(18)
	if err != nil {
		return StatusError{500, err, nil}
	}
	// Add prefix token is id of account
	randomToken = strconv.Itoa(user.ID) + randomToken

	// FLow 2: Create a new request forgot password
	txHandle, err := e.DB.Begin()
	if err != nil {
		return StatusError{500, err, nil}
	}
	pwChangeReq := models.NewReqForgotPw(user.ID, randomToken)
	err = pwChangeReq.Add(txHandle)
	if err != nil {
		logrus.Info("POST /user/forgot-password > Add ReqForgotPw")
		return StatusError{500, err, nil}
	}
	if pwChangeReq.ID == -1 {
		return StatusError{201, nil, []byte(`{"error": "Something went wrong. please try again"}`)}
	}

	// FLow 3: Create successful. Send a email for user and attach url token
	err = services.SendEmailForgotPassword(user.Email, pwChangeReq.Token)
	if err != nil {
		return StatusError{500, err, nil}
	}
	txHandle.Commit()
	return StatusError{200, nil, []byte(`{"error": ""}`)}
}

// ResetPassword swagger:route POST /user/reset-password Users ResetPassword
//
// 	Lấy tất cả những thông tin cơ bản của user
//
// 	Responses:
//		500: DocResDefault
func (e *Env) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Info("POST /user/reset-password > ResetPassword > ReadAll Error")
		return StatusError{500, err, nil}
	}
	var dataReq struct {
		NewPassword  string `json:"new_pw"`
		ConfirmNewPw string `json:"confirm_new_pw"`
		Token        string `json:"token"`
	}
	err = json.Unmarshal(b, &dataReq)
	if err != nil {
		logrus.Info("POST /user/reset-password > ResetPassword > Unmarshal Error")
		return StatusError{500, err, nil}
	}
	if dataReq.NewPassword != dataReq.ConfirmNewPw {
		logrus.Info("POST /user/reset-password > ResetPassword > Error")
		return StatusError{201, err, []byte(`{"error": "Password not match"}`)}
	}
	txHandle, err := e.DB.Begin()
	reqForgotPw, err := models.FindReqForgotPwByToken(txHandle, dataReq.Token)
	if err != nil {
		logrus.Info("POST /user/reset-password > ResetPassword > FindReqForgotPwByToken Error")
		return StatusError{500, err, nil}
	}

	acc, err := models.FindUserByID(e.DB, reqForgotPw.UID)
	if err != nil {
		logrus.Info("POST /user/reset-password > ResetPassword > FindUserByID Error")
		return StatusError{500, err, nil}
	}
	err = acc.ResetPassword(txHandle, dataReq.NewPassword)
	if err != nil {
		logrus.Info("POST /user/reset-password > ResetPassword Error")
		return StatusError{500, err, nil}
	}
	txHandle.Commit()
	return StatusError{200, nil, []byte(`{"message": "Your password has been changed successfully!"}`)}
}

// GetBalance ...
func (e *Env) GetBalance(w http.ResponseWriter, r *http.Request) error {
	userOauth, err := GetUserByOauthContext(r)
	if err != nil {
		logrus.Info("GET /userinfo > GetUserByOauthContext")
		return StatusError{201, err, []byte(`{"error": "Unauthorized"}`)}
	}

	userID := userOauth.ID
	moneys, err := models.SumBalanceByUserID(e.DB, userID)
	if err != nil {
		return err
	}
	resByteMoneys, err := json.Marshal(moneys)
	if err != nil {
		logrus.Info("GET /userinfo > GetUserInfo Error")
		return StatusError{500, err, nil}
	}
	return StatusError{200, nil, resByteMoneys}
}
