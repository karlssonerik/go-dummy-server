package logger

import "go.uber.org/zap"

var logger *zap.Logger
var sugar *zap.SugaredLogger

func init() {
	logger, _ = zap.NewProduction()
	sugar = logger.Sugar()
}

func Log() *zap.SugaredLogger {
	return sugar
}
