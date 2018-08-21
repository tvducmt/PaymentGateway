// Package classification Payment Gateway API.
//
// Description ... Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book.
//
//     Schemes: http, https
//     Host: localhost
//     Version: 0.1.0
//     License:
//     Contact: https://rockship.co/
//
//     Consumes:
//     - application/json
//     - application/xml
//
//     Produces:
//     - application/json
//     - application/xml
// swagger:meta
package api

import (
	stripe "github.com/stripe/stripe-go"
	"gitlab.com/rockship/payment-gateway/models"
)

//swagger:parameters Login
type DocAccountAuth struct {
	// in:body
	// required: true
	Body AccountAuth
}

//swagger:parameters Register
type DocAccRegister struct {
	// in: body
	Body struct {
		// required: true
		Email string `json:"email"`
		// required: true
		Password string `json:"password"`
		// required: true
		PwConfirm string `json:"pwConfirm"`
	}
}

//swagger:parameters UpdateUserInfo
type DocUpdateUserInfo struct {
	// in: body
	Body models.UserInfo
}

//swagger:parameters UpdatePassword
type DocUpdatePassword struct {
	// in: body
	Body models.Password
}

//swagger:parameters GetUserInfo GetCouponInfo
type DocUserID struct {
	// in:path
	UserID int `json:"userid"`
}

//swagger:parameters HistoryTransaction
type DocUserIDBody struct {
	// in:body
	Body struct {
		// min: 1
		// required: true
		UserID int `json:"userid"`
	}
}

//swagger:parameters GetTransactionDetail
type DocTxID struct {
	// in:path
	// required: true
	TxID int `json:"tx_id"`
}

//swagger:parameters ReceiveReceipt
type DocReceiveReceipt struct {
	// in:body
	// required: true
	Body struct {
		// required: true
		Timestamp string `json:"timestamp"`
		// required: true
		RawReceipt string `json:"receipt"`
		// required: true
		Phonenumber string `json:"phonenumber"`
	}
}

//swagger:parameters CreditCardProcess
type DocCreditCard struct {
	// in:body
	Body struct {
		// Thông tin thêm
		// required: false
		Name string `json:"name"`

		// Số tiền thực hiện giao dịch, mặc định đơn vị là vnd.
		// required: true
		Amount int64 `json:"amount"`

		// Token xác thực user
		// required: true
		Token string `json:"token"`

		// Lựa chọn của khách hàng sử dụng thẻ cũ hoặc dùng thẻ mới để giao dịch.
		// required: true
		ModeCard string `json:"select_mode_card"`

		// Lựa chọn một trong danh sách các thẻ cũ khách hàng đã lưu trước đó.
		// required: true
		CusID int `json:"select_cus_id"`

		// Yêu cầu có lưu lại thông tin thẻ hay không
		// required: true
		Remember bool `json:"remember_me"`
	}
}

//swagger:parameters BankTransferProcess
type DocBankTransfer struct {
	// in: body
	Body struct {
		// required: true
		UserID int `json:"user_id"`
		// required: true
		MethodID int `json:"method_id"`
		// required: true
		GatewayID int `json:"gateway_id"`
	}
}

// Server xử lí yêu cầu thành công
// swagger:response DocReceiveReceiptSuccess
type DocReceiveReceiptSuccess struct {
	// in: body
	Body struct {
		Error string `json:"error"`
		// required: true
		Coupon *models.Coupon `json:"coupon"`
	}
}

// Server xử lí yêu cầu thành công
// swagger:response DocResRegisterSuccess
type DocResRegisterSuccess struct {
	// in: body
	Body struct {
		// required: true
		Token string `json:"token"`
		Error string `json:"error"`
	}
}

// Update thành công
// swagger:response DocResUpdateUserInfo
type DocResUpdateUserInfo struct {
	//in: body
	Body struct {
		Message string `json:"message"`
	}
}

// Update thành công
// swagger:response DocResUpdateTxTimeOut
type DocResUpdateTxTimeOut struct {
	//in: body
}

// Server xử lí yêu cầu thành công
// swagger:response DocResUserInfo
type DocResUserInfo struct {
	// in: body
	Body models.UserInfo
}

// Server xử lí thành công. Lỗi phía người dùng
// swagger:response DocResUserMsg
type DocResUserMsg struct {
	// in: body
	Body struct {
		// UserID Invalid | User not exist
		Error string `json:"error"`
	}
}

// Server xử lí thành công. Lỗi phía người dùng
// swagger:response DocResTxMsg
type DocResTxMsg struct {
	// in: body
	Body struct {
		// Password Not Match | Existing Account | Error establishing a database connection
		Error string `json:"error"`
	}
}

// Server xử lí thành công. Lỗi phía người dùng
// swagger:response DocResRegister
type DocResRegister struct {
	// in: body
	Body struct {
		// Password Not Match | Existing Account | Error establishing a database connection
		Error string `json:"error"`
	}
}

// Lỗi phía Server
// swagger:response DocResDefault
type DocResDefault struct {
	// in: body
	Body struct {
		// Something went wrong. please try again
		Error string `json:"error"`
	}
}

// Server xử lí thành công
// swagger:response DocResCouponInfo
type DocResCouponInfo struct {
	// in: body
	Body struct {
		Error   string          `json:"error"`
		Coupons []models.Coupon `json:"coupons"`
	}
}

// Server xử lí thành công
// swagger:response DecResCustomer
type DecResCustomer struct {
	// in: body
	Body struct {
		CusUsers []struct {
			// ID định danh cho customer
			//
			// min: 1
			ID int `json:"id"`

			// UserID khóa ngoại đến model user
			//
			// min: 1
			UserID int `json:"user_id"`

			// Bốn số cuối thẻ
			Last4 string `json:"last4"`

			// Tháng hết hạn
			ExpMonth uint8 `json:"exp_month"`

			// Năm hết hạn
			ExpYear uint16 `json:"exp_year"`

			// Mã bưu điện
			AddressZip string `json:"address_zip"`

			// Loại thẻ (Visa, MasterCard, ...)
			Brand string `json:"brand"`

			// Hình thức (Credit, Debit, ...)
			Funding string `json:"funding"`

			// Chữ ký trên thẻ
			// Unique: true
			Fingerprint string `json:"fingerprint"`
		} `json:"data"`
	}
}

// Server xử lí thành công
// swagger:response DocResTxSuccess
type DocResTxSuccess struct {
	// in: body
	Body struct {
		Transaction models.TransactionUser `json:"transaction"`
		MetaCredit  *stripe.Charge         `json:"charge"`
	}
}

// Server xử lí thành công
// swagger:response DocResLogin
type DocResLogin struct {
	// in: body
	Body struct {
		// OK | ERROR
		Status string `json:"status"`

		// '' | Email or Password wrong!
		Error string `json:"error"`

		// required: true
		Token string `json:"token"`
	}
}
