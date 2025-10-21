package logger

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func New() *Logger {

	return &Logger{
		Logger: log.New(os.Stdout, "[3xui-bot] ", log.LstdFlags|log.Lshortfile),
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.Printf("[INFO] "+format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.Printf("[ERROR] "+format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.Printf("[DEBUG] "+format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.Fatalf("[FATAL] "+format, v...)
}
