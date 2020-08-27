package log

import (
	"go.uber.org/zap"
)

var SugarLogger *zap.SugaredLogger
var Logger *zap.Logger

func Init() {
	Logger, _ = zap.NewProduction()
	SugarLogger = Logger.Sugar()
}

func Stop() {
	Logger.Sync()
	SugarLogger.Sync()
}
