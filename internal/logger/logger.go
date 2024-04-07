package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

var log *Logger

func InitSlog(level string) {
	var slogLevel slog.Leveler
	switch level {
	case "INFO":
		slogLevel = slog.LevelInfo
	case "DEBUG":
		slogLevel = slog.LevelDebug
	case "ERROR":
		slogLevel = slog.LevelError
	}
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel}))
	log = &Logger{logger: slogger}
}

func (l *Logger) Info(title string, msg ...any) {
	l.logger.Info(title, msg...)
}

func (l *Logger) Warn(title string, msg ...any) {
	l.logger.Warn(title, msg...)
}

func (l *Logger) Debug(title string, msg ...any) {
	l.logger.Debug(title, msg...)
}
func (l *Logger) Error(title string, err ...any) {
	l.logger.Error(title, err...)
}

func Log() *Logger {
	return log
}
