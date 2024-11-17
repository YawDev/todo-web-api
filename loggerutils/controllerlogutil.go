package loggerutils

import (
	"context"

	"github.com/sirupsen/logrus"
)

func ErrorLog(ctx context.Context, statusCode int, err error) {
	Log.FromContext(ctx).WithFields(logrus.Fields{
		"Status": statusCode,
	}).Error(err.Error())
}

func InfoLog(ctx context.Context, statusCode int, msg string) {
	Log.FromContext(ctx).WithFields(logrus.Fields{
		"Status": statusCode,
	}).Info(msg)

}
