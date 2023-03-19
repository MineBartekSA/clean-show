package logger

import "go.uber.org/zap"

var Log *zap.SugaredLogger

func InitProduction() {
	logger, _ := zap.NewProduction()
	Log = logger.Sugar()
}

func InitDebug() {
	logger, _ := zap.NewDevelopment()
	Log = logger.Sugar()
}
