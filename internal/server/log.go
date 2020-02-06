package server

import "fmt"

type LogLevel int

var logLevel = LogLevelInfo

func SetLogLevel(lvl LogLevel) {
	logLevel = lvl
}

const (
	LogLevelDebug = LogLevel(0)
	LogLevelInfo  = LogLevel(1)
	LogLevelError = LogLevel(2)
)

func logDebug(msg string) {
	if logLevel <= LogLevelDebug {
		fmt.Println("[DEBUG] " + msg)
	}
}

func logInfo(msg string) {
	if logLevel <= LogLevelInfo {
		fmt.Println("[INFO] " + msg)
	}
}
func logError(msg string) {
	if logLevel <= LogLevelError {
		fmt.Println("[ERROR] " + msg)
	}
}
