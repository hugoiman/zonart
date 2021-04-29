package controllers

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

// Gomail is class
type Gomail struct {
	configSMTPHost string
	configSMTPPort int
	configEmail    string
	configPassword string
}

func (g *Gomail) GetVariable() Gomail {
	g.configSMTPHost = "smtp.gmail.com"
	g.configSMTPPort = 587
	g.configEmail = "nanonymoux@gmail.com"
	g.configPassword = "bkl"

	return *g
}

// SendEmail is func
func (g Gomail) SendEmail(subject string, to string, message string) error {
	env := g.GetVariable()
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", env.configEmail)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", message)

	dialer := gomail.NewDialer(env.configSMTPHost, env.configSMTPPort, env.configEmail, env.configPassword)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := dialer.DialAndSend(mailer)

	return err

}
