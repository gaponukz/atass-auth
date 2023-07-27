package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Settings struct {
	JwtSecret        string `json:"jwtSecret"`
	HashSecret       string `json:"hashSecret"`
	Gmail            string `json:"gmail"`
	GmailPassword    string `json:"gmailPassword"`
	Port             int64  `json:"port"`
	RedisAddress     string `json:"redisAddress"`
	PostgresHost     string `json:"postgresHost"`
	PostgresUser     string `json:"postgresUser"`
	PostgresPassword string `json:"postgresPassword"`
	PostgresDbname   string `json:"postgresDbname"`
	PostgresPort     string `json:"postgresPort"`
	PostgresSslmode  string `json:"postgresSslmode"`
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

func checkRedis(r string) string {
	if r == "" {
		return "localhost:6379"
	}

	return r
}

func (sts dotEnvSettings) Load() Settings {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
	}

	return Settings{
		JwtSecret:        os.Getenv("jwtSecret"),
		HashSecret:       os.Getenv("hashSecret"),
		Gmail:            os.Getenv("gmail"),
		GmailPassword:    os.Getenv("gmailPassword"),
		Port:             parsePort(os.Getenv("port")),
		RedisAddress:     checkRedis(os.Getenv("redisAddress")),
		PostgresHost:     os.Getenv("postgresHost"),
		PostgresUser:     os.Getenv("postgresUser"),
		PostgresPassword: os.Getenv("postgresPassword"),
		PostgresDbname:   os.Getenv("postgresDbname"),
		PostgresPort:     os.Getenv("postgresPort"),
		PostgresSslmode:  os.Getenv("postgresSslmode"),
	}
}
