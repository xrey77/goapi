package users

import (
	"crypto/tls"

	"gopkg.in/mail.v2"
)

func ActivateAccount(msg string, sub string, xmail string) int64 {
	m := mail.NewMessage()
	m.SetHeader("From", "rey107@gmail.com")
	m.SetHeader("To", xmail)
	m.SetHeader("Subject", sub)
	m.SetBody("text/html", msg)
	//m.Attach("/home/Alex/lolcat.jpg")
	// Settings for SMTP server
	d := mail.NewDialer("smtp.gmail.com", 587, "rey107@gmail.com", "Reynald@88.88")
	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		return 0
	}
	return 1
}
