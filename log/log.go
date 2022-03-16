package log

import (
	"go.uber.org/zap"
)

var SugarLogger *zap.SugaredLogger
var Logger *zap.Logger

func init() {
	Logger, _ = zap.NewProduction()
	SugarLogger = Logger.Sugar()
}
