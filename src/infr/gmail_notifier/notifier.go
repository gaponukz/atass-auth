package gmail_notifier

import (
	"auth/src/domain/entities"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

type GmailCreds struct {
	Gmail    string
	Password string
}

type Letter struct {
	Title    string
	HtmlPath string
}

type gmailNotifier struct {
	auth       smtp.Auth
	gmail      string
	letterHtml string
	title      string
}

func NewGmailNotifier(creds GmailCreds, letter Letter) *gmailNotifier {
	if !isValidTemplatePath(letter.HtmlPath) {
		return nil
	}

	data, err := os.ReadFile(filepath.Join("", filepath.Clean(letter.HtmlPath)))
	if err != nil {
		return nil
	}

	return &gmailNotifier{
		auth:       smtp.PlainAuth("", creds.Gmail, creds.Password, "smtp.gmail.com"),
		letterHtml: string(data),
		gmail:      creds.Gmail,
		title:      letter.Title,
	}
}

func (n gmailNotifier) Notify(to, code string) error {
	return smtp.SendMail("smtp.gmail.com:587", n.auth, n.gmail, []string{to}, n.generateLetter(to, code))
}

func (n gmailNotifier) NotifyUser(to entities.User, code string) error {
	return n.Notify(to.Gmail, code)
}

func (n gmailNotifier) generateLetter(to, code string) []byte {
	body := strings.Replace(n.letterHtml, "CONFIRMATION_CODE", code, -1)
	return []byte(
		"To: " + to + "\r\n" +
			"Subject: " + n.title + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			body + "\r\n",
	)
}

func isValidTemplatePath(path string) bool {
	return strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".htm")
}
