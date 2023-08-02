package notifier

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

type Notifier func(to string, code string) error

type SendFrom struct {
	Gmail    string
	Password string
}

func SendEmailNoificationFactory(sender SendFrom, title string, templatePath string) Notifier {
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil
	}

	return func(sendToGmail string, code string) error {
		body := strings.Replace(string(data), "CONFIRMATION_CODE", code, -1)
		message := []byte(
			"To: " + sendToGmail + "\r\n" +
				"Subject: " + title + "\r\n" +
				"MIME-Version: 1.0\r\n" +
				"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
				"\r\n" +
				body + "\r\n",
		)

		auth := smtp.PlainAuth("", sender.Gmail, sender.Password, "smtp.gmail.com")

		err := smtp.SendMail("smtp.gmail.com:587", auth, sender.Gmail, []string{sendToGmail}, message)
		if err != nil {
			return fmt.Errorf("Can not send letter: %v", err)
		}

		return nil
	}
}
