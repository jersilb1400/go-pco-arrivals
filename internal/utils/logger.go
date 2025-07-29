package utils

import (
	"log"
	"os"
	"time"
)

type Logger struct {
	component string
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{
		component: component,
	}
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.log("INFO", message, args...)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.log("WARN", message, args...)
}

func (l *Logger) Error(message string, args ...interface{}) {
	l.log("ERROR", message, args...)
}

func (l *Logger) Debug(message string, args ...interface{}) {
	l.log("DEBUG", message, args...)
}

func (l *Logger) Fatal(message string, args ...interface{}) {
	l.log("FATAL", message, args...)
	os.Exit(1)
}

func (l *Logger) log(level, message string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	if l.component != "" {
		log.Printf("[%s] [%s] [%s] %s", timestamp, level, l.component, message)
	} else {
		log.Printf("[%s] [%s] %s", timestamp, level, message)
	}

	// Log additional arguments if provided
	if len(args) > 0 {
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				log.Printf("  %v: %v", args[i], args[i+1])
			}
		}
	}
}
