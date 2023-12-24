package logger

import (
	"go.uber.org/zap"
	"testing"
)

func InitTestLogger(t *testing.T) {
	var err error
	Main, err = zap.NewDevelopment()
	if err != nil {
		t.Error(err.Error())
	}
}
