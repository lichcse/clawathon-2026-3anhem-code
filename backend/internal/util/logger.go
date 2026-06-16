package util

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
	LevelDebug Level = "DEBUG"
)

type Logger struct {
	out    io.Writer
	fields map[string]any
}

type logEntry struct {
	Time    string         `json:"time"`
	Level   Level          `json:"level"`
	Message string         `json:"message"`
	Fields  map[string]any `json:"fields,omitempty"`
}

var defaultLogger = &Logger{out: os.Stdout, fields: map[string]any{}}

func NewLogger(out io.Writer) *Logger {
	return &Logger{out: out, fields: map[string]any{}}
}

func (l *Logger) With(key string, value any) *Logger {
	newFields := make(map[string]any, len(l.fields)+1)
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = value
	return &Logger{out: l.out, fields: newFields}
}

func (l *Logger) log(level Level, msg string, extra map[string]any) {
	entry := logEntry{
		Time:    time.Now().UTC().Format(time.RFC3339),
		Level:   level,
		Message: msg,
	}
	if len(l.fields) > 0 || len(extra) > 0 {
		fields := make(map[string]any, len(l.fields)+len(extra))
		for k, v := range l.fields {
			fields[k] = v
		}
		for k, v := range extra {
			fields[k] = v
		}
		entry.Fields = fields
	}
	b, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger marshal error: %v\n", err)
		return
	}
	fmt.Fprintln(l.out, string(b))
}

func (l *Logger) Info(msg string, fields ...map[string]any) {
	extra := map[string]any{}
	if len(fields) > 0 {
		extra = fields[0]
	}
	l.log(LevelInfo, msg, extra)
}

func (l *Logger) Infof(format string, args ...any) {
	l.log(LevelInfo, fmt.Sprintf(format, args...), nil)
}

func (l *Logger) Warn(msg string, fields ...map[string]any) {
	extra := map[string]any{}
	if len(fields) > 0 {
		extra = fields[0]
	}
	l.log(LevelWarn, msg, extra)
}

func (l *Logger) Warnf(format string, args ...any) {
	l.log(LevelWarn, fmt.Sprintf(format, args...), nil)
}

func (l *Logger) Error(msg string, fields ...map[string]any) {
	extra := map[string]any{}
	if len(fields) > 0 {
		extra = fields[0]
	}
	l.log(LevelError, msg, extra)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.log(LevelError, fmt.Sprintf(format, args...), nil)
}

func (l *Logger) Debug(msg string, fields ...map[string]any) {
	extra := map[string]any{}
	if len(fields) > 0 {
		extra = fields[0]
	}
	l.log(LevelDebug, msg, extra)
}

func (l *Logger) WithTime(format string, args ...any) {
	l.log(LevelInfo, fmt.Sprintf(format, args...), nil)
}

// Package-level helpers using defaultLogger
func Info(msg string)           { defaultLogger.Info(msg) }
func Infof(f string, a ...any)  { defaultLogger.Infof(f, a...) }
func Warn(msg string)           { defaultLogger.Warn(msg) }
func Warnf(f string, a ...any)  { defaultLogger.Warnf(f, a...) }
func Error(msg string)          { defaultLogger.Error(msg) }
func Errorf(f string, a ...any) { defaultLogger.Errorf(f, a...) }
func Debug(msg string)          { defaultLogger.Debug(msg) }
func SetLogger(l *Logger)       { defaultLogger = l }
