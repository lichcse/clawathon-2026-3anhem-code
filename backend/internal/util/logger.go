package util

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	warnLog  *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		infoLog:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLog:  log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime),
	}
}

func (l *Logger) Info(msg string) {
	l.infoLog.Println(msg)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.infoLog.Println(fmt.Sprintf(format, args...))
}

func (l *Logger) Error(msg string) {
	l.errorLog.Println(msg)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.errorLog.Println(fmt.Sprintf(format, args...))
}

func (l *Logger) Warn(msg string) {
	l.warnLog.Println(msg)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.warnLog.Println(fmt.Sprintf(format, args...))
}

func (l *Logger) WithTime(msg string) string {
	return fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), msg)
}
