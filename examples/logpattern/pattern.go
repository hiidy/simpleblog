package main

import (
	"fmt"
	"strings"
)

const (
	DebugLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type customLogger struct {
	level  int
	prefix string
}

// Debug implements [Logger].
func (l *customLogger) Debug(msg string) {
	l.logMessage(DebugLevel, "debug", msg)
}

// Error implements [Logger].
func (l *customLogger) Error(msg string) {
	l.logMessage(InfoLevel, "info", msg)
}

// Info implements [Logger].
func (l *customLogger) Info(msg string) {
	l.logMessage(WarnLevel, "warn", msg)
}

// Warn implements [Logger].
func (l *customLogger) Warn(msg string) {
	l.logMessage(ErrorLevel, "error", msg)
}

func New(level int, prefix string) Logger {
	return &customLogger{
		level:  level,
		prefix: prefix,
	}
}

func (l *customLogger) logMessage(level int, levelStr string, msg string) {
	if level >= l.level {
		fmt.Printf("[%s] %s: %s\n", strings.ToUpper(levelStr), l.prefix)
	}
}

func main() {
	logger := New(InfoLevel, "MyApp")

	logger.Debug("This is a debug message")
	logger.Info("THis is an info message")
}
