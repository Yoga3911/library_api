package helper

import (
	"log"

	gomail "gopkg.in/mail.v2"
)

func SendOTP(otp string, receiver string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "yoga3911@gmail.com")
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", "OTP Code")
	m.SetBody("text/plain", otp)
	m.Attach("./assets/images/bagan.png")

	d := gomail.NewDialer("smtp.gmail.com", 587, "yoga3911@gmail.com", "zjwwbwyxwucovmvj")

	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
	}
}
