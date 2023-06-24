package settings

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Settings struct {
	JwtSecret     string `json:"jwtSecret"`
	Gmail         string `json:"gmail"`
	GmailPassword string `json:"gmailPassword"`
	Port          int64  `json:"port"`
}

type dotEnvSettings struct{}

func NewDotEnvSettings() *dotEnvSettings {
	return &dotEnvSettings{}
}

func parsePort(port string) int64 {
	i, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		return 8000
	}

	return i
}

func (sts dotEnvSettings) Load() Settings {
	err := godotenv.Load()
	if err != nil {
		return Settings{}
	}

	return Settings{
		JwtSecret:     os.Getenv("jwtSecret"),
		Gmail:         os.Getenv("gmail"),
		GmailPassword: os.Getenv("gmailPassword"),
		Port:          parsePort(os.Getenv("port")),
	}
}
