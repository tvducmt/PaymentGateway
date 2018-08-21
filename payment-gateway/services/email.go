package services

import (
	"bytes"
	"html/template"
	"net/smtp"
	"os"
)

func SendEmailForgotPassword(email, token string) error {
	users := []string{email}
	auth := smtp.PlainAuth("", os.Getenv("EMAIL"), os.Getenv("PW_EMAIL"), "smtp.gmail.com")
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	mime := " \nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: Quên mật khẩu Rockship Payment Gateway API"

	body, err := emailTemplateHTML(token)
	if err != nil {
		return err
	}

	err = smtp.SendMail("smtp.gmail.com:25",
		auth, os.Getenv("EMAIL"), users,
		[]byte(subject+mime+"\n"+body))
	if err != nil {
		return err
	}
	return nil
}

func emailTemplateHTML(token string) (string, error) {
	t, err := template.ParseFiles(os.Getenv("TEMPLATE_HTML"))
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	url := os.Getenv("BASE_URL_FORGOT_PASSWORD_UI") + token
	if err := t.Execute(&tpl, url); err != nil {
		return "", err
	}

	result := tpl.String()
	return result, nil
}
