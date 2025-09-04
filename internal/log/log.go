package log

import (
	"fmt"
	"os"
	"time"
)

// Couleurs ANSI
const (
	reset   = "\033[0m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"
)

var Logs = []string{}

func Info(args ...interface{}) {
	go log("INFO", cyan, args...)
}

func Success(args ...interface{}) {
	go log("SUCCESS", green, args...)
}

func Warn(args ...interface{}) {
	go log("WARN", yellow, args...)
}

func Error(args ...interface{}) {
	go log("ERROR", red, args...)
}

func Debug(args ...interface{}) {
	go log("DEBUG", magenta, args...)
}

func Fatal(message string, err error) {
	fmt.Printf("FATAL: %s: %s\n", message, err)
	os.Exit(1)
}

func log(level string, color string, args ...interface{}) {
	now := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprint(args...)

	levelBadge := fmt.Sprintf("%s%s%s", color, level, reset)

	Logs = append(Logs, fmt.Sprintf("[%s] %s %s", now, levelBadge, msg))
}

func WriteLogs() {
	for _, log := range Logs {
		fmt.Println(log)
	}

	Logs = []string{}
}
