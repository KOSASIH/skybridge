package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// Logger represents a custom logger type
type Logger struct {
	level  LogLevel
	output io.Writer
	mu     sync.Mutex
}

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// NewLogger creates a new logger instance
func NewLogger(level LogLevel, output io.Writer) *Logger {
	return &Logger{
		level:  level,
		output: output,
	}
}

// Log logs a message at the specified level
func (l *Logger) Log(level LogLevel, message string, args ...interface{}) {
	if level < l.level {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	log.SetOutput(l.output)
	log.Printf("[%s] %s", levelToString(level), fmt.Sprintf(message, args...))
}

// Debug logs a message at the DEBUG level
func (l *Logger) Debug(message string, args ...interface{}) {
	l.Log(DEBUG, message, args...)
}

// Info logs a message at the INFO level
func (l *Logger) Info(message string, args ...interface{}) {
	l.Log(INFO, message, args...)
}

// Warning logs a message at the WARNING level
func (l *Logger) Warning(message string, args ...interface{}) {
	l.Log(WARNING, message, args...)
}

// Error logs a message at the ERROR level
func (l *Logger) Error(message string, args ...interface{}) {
	l.Log(ERROR, message, args...)
}

// Fatal logs a message at the FATAL level and exits the program
func (l *Logger) Fatal(message string, args ...interface{}) {
	l.Log(FATAL, message, args...)
	os.Exit(1)
}

func levelToString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}
