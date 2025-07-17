package logger

import (
	"github.com/charmbracelet/log"
)

// Hijack the wails' logger
type WailsLogger struct{}

func (l *WailsLogger) Print(msg string) {
	log.Print(msg)
}

func (l *WailsLogger) Trace(msg string) {
	log.Debug(msg)
}

func (l *WailsLogger) Debug(msg string) {
	log.Debug(msg)
}

func (l *WailsLogger) Info(msg string) {
	log.Info(msg)
}

func (l *WailsLogger) Warning(msg string) {
	log.Warn(msg)
}

func (l *WailsLogger) Error(msg string) {
	log.Error(msg)
}

func (l *WailsLogger) Fatal(msg string) {
	log.Error(msg)
}
