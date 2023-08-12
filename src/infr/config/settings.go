package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Settings struct {
	JwtSecret     string `json:"jwtSecret"`
	HashSecret    string `json:"hashSecret"`
	Gmail         string `json:"gmail"`
	GmailPassword string `json:"gmailPassword"`
	Port          int64  `json:"port"`
	RedisAddress  string `json:"redisAddress"`
	MysqlHost     string `json:"mysqlHost"`
	MysqlUser     string `json:"mysqlUser"`
	MysqlPassword string `json:"mysqlPassword"`
	MysqlDbname   string `json:"mysqlDbname"`
	MysqlPort     string `json:"mysqlPort"`
	MysqlSslmode  string `json:"mysqlSslmode"`
	RabbitUrl     string `json:"rabbitUrl"`
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
		JwtSecret:     os.Getenv("jwtSecret"),
		HashSecret:    os.Getenv("hashSecret"),
		Gmail:         os.Getenv("gmail"),
		GmailPassword: os.Getenv("gmailPassword"),
		Port:          parsePort(os.Getenv("port")),
		RedisAddress:  checkRedis(os.Getenv("redisAddress")),
		MysqlHost:     os.Getenv("mysqlHost"),
		MysqlUser:     os.Getenv("mysqlUser"),
		MysqlPassword: os.Getenv("mysqlPassword"),
		MysqlDbname:   os.Getenv("mysqlDbname"),
		MysqlPort:     os.Getenv("mysqlPort"),
		MysqlSslmode:  os.Getenv("mysqlSslmode"),
		RabbitUrl:     os.Getenv("rabbitmqUrl"),
	}
}
