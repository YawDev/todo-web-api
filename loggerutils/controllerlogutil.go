package loggerutils

import (
	"context"
)

func ErrorLog(ctx context.Context, statusCode int, err error) {
	Log.FromContext(ctx).Error(err.Error())
}

func InfoLog(ctx context.Context, statusCode int, msg string) {
	Log.FromContext(ctx).Info(msg)

}
