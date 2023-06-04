package settings

import (
	"os"

	"github.com/joho/godotenv"
)

type Settings struct {
	JwtSecret     string `json:"jwtSecret"`
	Gmail         string `json:"gmail"`
	GmailPassword string `json:"gmailPassword"`
}

type DotEnvSettings struct{}

func (sts DotEnvSettings) Load() Settings {
	godotenv.Load()

	return Settings{
		JwtSecret:     os.Getenv("jwtSecret"),
		Gmail:         os.Getenv("gmail"),
		GmailPassword: os.Getenv("gmailPassword"),
	}
}
