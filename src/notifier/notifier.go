package notifier

import (
	"net/smtp"
)

type SendFrom struct {
	Gmail    string
	Password string
}

func SendEmailNoificationFactory(sender SendFrom) func(to string, title string, body string) error {
	return func(sendToGmail string, title string, body string) error {
		message := []byte(
			"To: " + sendToGmail + "\r\n" +
				"Subject: " + title + "\r\n" +
				"MIME-Version: 1.0\r\n" +
				"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
				"\r\n" +
				body + "\r\n",
		)

		auth := smtp.PlainAuth("", sender.Gmail, sender.Password, "smtp.gmail.com")
		return smtp.SendMail("smtp.gmail.com:587", auth, sender.Gmail, []string{sendToGmail}, message)
	}
}
