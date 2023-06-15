package notifier

import (
	"os"
	"strings"
)

func GenerateConfirmCodeLetter(path, code string) string {
	data, err := os.ReadFile(path)

	if err != nil {
		panic("file not found")
	}

	return strings.Replace(string(data), "CONFIRMATION_CODE", code, -1)
}
