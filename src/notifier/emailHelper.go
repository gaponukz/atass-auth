package notifier

import (
	"os"
	"strings"
)

func GenerateConfirmCodeLetter(code string) string {
	data, err := os.ReadFile("letters/letter.html")

	if err != nil {
		panic("file not found")
	}

	return strings.Replace(string(data), "CONFIRMATION_CODE", code, -1)
}
