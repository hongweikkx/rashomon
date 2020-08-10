package log

import (
	"go.uber.org/zap"
)

var SugarLogger *zap.SugaredLogger
var Logger *zap.Logger

func InitLogger() {
	Logger, _ = zap.NewProduction()
	SugarLogger = Logger.Sugar()
}
