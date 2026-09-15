package utils

import "go.uber.org/zap"

var Logger *zap.SugaredLogger

func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}

	Logger = logger.Sugar()
}
