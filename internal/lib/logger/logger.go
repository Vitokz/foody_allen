package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func NewLogger(mode string) (*zap.SugaredLogger, error) {
	var l *zap.Logger
	var err error

	if mode == "development" {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}
	if err != nil {
		return nil, err
	}

	Logger = l.Sugar()

	return Logger, nil
}
