package logger

import "github.com/charmbracelet/log"

// A simple wrapper of charmbracelet/log, to hijack the wails' logger
type Logger struct {
}

func (l *Logger) Print(msg string) {
	log.Print(msg)
}

func (l *Logger) Trace(msg string) {
	log.Debug(msg)
}

func (l *Logger) Debug(msg string) {
	log.Debug(msg)
}

func (l *Logger) Info(msg string) {
	log.Info(msg)
}

func (l *Logger) Warning(msg string) {
	log.Warn(msg)
}

func (l *Logger) Error(msg string) {
	log.Error(msg)
}

func (l *Logger) Fatal(msg string) {
	log.Error(msg)
}
