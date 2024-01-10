package logging

import (
	"fmt"
	"log"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

func Info(format string, a ...interface{}) {
	log.Printf(InfoColor, fmt.Sprintf(format, a...))
}

func Notice(format string, a ...interface{}) {
	log.Printf(NoticeColor, fmt.Sprintf(format, a...))
}

func Warning(format string, a ...interface{}) {
	log.Printf(WarningColor, fmt.Sprintf(format, a...))
}

func Error(format string, a ...interface{}) {
	log.Fatalf(ErrorColor, fmt.Sprintf(format, a...))
}

func Debug(format string, a ...interface{}) {
	log.Printf(DebugColor, fmt.Sprintf(format, a...))
}
