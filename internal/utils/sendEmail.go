package utils

import (
	"log"
	"net/smtp"

	"github.com/Dimpal-Kalita/LetsNormalizeIt-backend/config"
)

// https://mailmeteor.com/blog/gmail-smtp-settings
// setting up only the app password in gmail works
// for reference

func SendVerificationEmail(to string, link string) error {
	cfg := config.Loadconfig()
	from := cfg.EMAIL_ID
	password := cfg.EMAIL_PASS
	BASE_URL := cfg.BASE_URL
	link = BASE_URL + "/verify/" + link
	msg := "From: " + from + "\n" + "To: " + to + "\n" + "Subject: Varification Email\n\n" + "Click on the link below to varify your email\n" + link

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, password, "smtp.gmail.com"), from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	log.Print("Verification email sent to: " + to)
	return nil
}

func SendResetPasswordEmail(to string, link string) error {
	cfg := config.Loadconfig()
	from := cfg.EMAIL_ID
	password := cfg.EMAIL_PASS
	BASE_URL := cfg.BASE_URL
	link = BASE_URL + "/reset-password/" + link
	msg := "From: " + from + "\n" + "To: " + to + "\n" + "Subject: Reset Password Email\n\n" + "Click on the link below to Reset your Password. The Link is valid till 60 minutes from now\n" + link

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, password, "smtp.gmail.com"), from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	log.Print("Password Reset Email sent to: " + to)
	return nil
}
