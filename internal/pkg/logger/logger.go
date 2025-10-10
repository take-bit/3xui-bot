package logger

import (
	"log"
	"os"
)

// Logger простой логгер (можно заменить на zap/zerolog)
type Logger struct {
	*log.Logger
}

// New создает новый логгер
func New() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "[3xui-bot] ", log.LstdFlags|log.Lshortfile),
	}
}

// Info логирует информационное сообщение
func (l *Logger) Info(format string, v ...interface{}) {
	l.Printf("[INFO] "+format, v...)
}

// Error логирует ошибку
func (l *Logger) Error(format string, v ...interface{}) {
	l.Printf("[ERROR] "+format, v...)
}

// Debug логирует отладочное сообщение
func (l *Logger) Debug(format string, v ...interface{}) {
	l.Printf("[DEBUG] "+format, v...)
}

// Fatal логирует фатальную ошибку и завершает программу
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.Fatalf("[FATAL] "+format, v...)
}
