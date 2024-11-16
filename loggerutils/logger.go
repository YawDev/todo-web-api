package loggerutils

import (
	"context"
	"os"
	"sync"
	"todo-web-api/contextkeys"

	"github.com/sirupsen/logrus"
)

var (
	once           sync.Once
	loggerInstance *LogUtil
)

type LogUtil struct {
	Logger *logrus.Logger
}

var Log = GetLogger()

func NewLogger() *LogUtil {
	l := logrus.New()

	l.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	l.SetLevel(logrus.InfoLevel)
	l.SetOutput(os.Stdout)

	return &LogUtil{Logger: l}
}

func GetLogger() *LogUtil {
	once.Do(func() {
		loggerInstance = NewLogger()
	})
	return loggerInstance
}

func (l *LogUtil) Info(msg string) {
	l.Logger.Info(msg)
}

func (l *LogUtil) Infof(msg string, args ...interface{}) {
	l.Logger.Infof(msg, args...)
}

func (l *LogUtil) Warn(msg string) {
	l.Logger.Warn(msg)
}

func (l *LogUtil) Error(msg string) {
	l.Logger.Error(msg)
}

func (l *LogUtil) Debug(msg string) {
	l.Logger.Debug(msg)
}

func (l *LogUtil) LogWithFields(msg string, fields logrus.Fields) {
	l.Logger.WithFields(fields).Info(msg)
}

func (l *LogUtil) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

func (l *LogUtil) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

func (l *LogUtil) Warningf(format string, args ...interface{}) {
	l.Logger.Warningf(format, args...)
}

func (l *LogUtil) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatalf(format, args...)
}

func (l *LogUtil) FromContext(ctx context.Context) *logrus.Entry {
	logger, ok := ctx.Value(contextkeys.ContextLoggerKey).(*logrus.Entry)
	if !ok {
		l.Logger.Warn("Logger not found in context")
		return logrus.NewEntry(logrus.StandardLogger())
	}
	l.Logger.Info("Logger found in context")
	return logger
}
