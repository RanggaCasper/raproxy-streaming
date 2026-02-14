package logger

import (
	"log"
	"os"
)

// Logger provides structured logging
type Logger struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

// New creates a new logger instance
func New() *Logger {
	return &Logger{
		errorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
	}
}

// Error logs error messages
func (l *Logger) Error(format string, v ...interface{}) {
	l.errorLog.Printf(format, v...)
}

// Info logs info messages
func (l *Logger) Info(format string, v ...interface{}) {
	l.infoLog.Printf(format, v...)
}
