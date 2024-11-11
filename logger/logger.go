package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	log *logrus.Logger
}

func NewLogger() *Logger {
	l := logrus.New()

	l.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	l.SetLevel(logrus.InfoLevel)
	l.SetOutput(os.Stdout)

	return &Logger{log: l}
}

func (l *Logger) Info(msg string) {
	l.log.Info(msg)
}

func (l *Logger) Warn(msg string) {
	l.log.Warn(msg)
}

func (l *Logger) Error(msg string) {
	l.log.Error(msg)
}

func (l *Logger) Debug(msg string) {
	l.log.Debug(msg)
}

func (l *Logger) LogWithFields(msg string, fields logrus.Fields) {
	l.log.WithFields(fields).Info(msg)
}

func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.log.WithFields(fields)
}
