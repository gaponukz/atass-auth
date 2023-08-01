package logger

import (
	"fmt"
	"time"
)

type consoleLogger struct{}

func NewConsoleLogger() consoleLogger {
	return consoleLogger{}
}

func (l consoleLogger) publishMessage(lvl, message string) {
	fmt.Println(l.formatMessage(lvl, message))
}

func (l consoleLogger) Info(message string) {
	l.publishMessage("INFO", message)
}

func (l consoleLogger) Warn(message string) {
	l.publishMessage("WARN", message)
}

func (l consoleLogger) Error(message string) {
	l.publishMessage("ERROR", message)
}

func (l consoleLogger) formatMessage(level, message string) string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%s | %s | %s", formattedTime, level, message)
}
