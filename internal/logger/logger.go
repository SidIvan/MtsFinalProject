package logger

import "go.uber.org/zap"

var Main *zap.Logger

func InitMainLogger(isDev bool) error {
	var err error
	if isDev {
		Main, err = zap.NewDevelopment()
		return err
	}
	Main, err = zap.NewProduction()
	return err
}
