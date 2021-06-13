package controllers

import (
	"crypto/tls"
	"os"

	"gopkg.in/gomail.v2"
)

// Gomail is class
type Gomail struct {
	configSMTPHost string
	configSMTPPort int
	configEmail    string
	configPassword string
}

func (g *Gomail) setVariable() {
	g.configSMTPHost = "smtp.gmail.com"
	g.configSMTPPort = 587
	g.configEmail = "nanonymoux@gmail.com"
	g.configPassword = os.Getenv("PASS_EMAIL")
}

// SendEmail is func
func (g Gomail) SendEmail(subject string, to string, message string) error {
	g.setVariable()
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", g.configEmail)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", message)

	dialer := gomail.NewDialer(g.configSMTPHost, g.configSMTPPort, g.configEmail, g.configPassword)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := dialer.DialAndSend(mailer)

	return err

}
