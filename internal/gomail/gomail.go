package gomail

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

// SendEmail is func
func SendEmail(subject string, to string, message string) error {
	var configSMTPHost = "smtp.gmail.com"
	var configSMTPPort = 587
	var configEmail = "nanonymoux@gmail.com"
	// var configPassword = os.Getenv("PASS_EMAIL")
	var configPassword = "bklHKshIS"

	mailer := gomail.NewMessage()

	mailer.SetHeader("From", configEmail)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", message)

	dialer := gomail.NewDialer(configSMTPHost, configSMTPPort, configEmail, configPassword)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := dialer.DialAndSend(mailer)

	return err

}
