package log

import (
	"os"

	"github.com/hongweikkx/rashomon/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger
var MLogger *zap.Logger

func init() {
	var logger *zap.Logger
	var err error
	if conf.AppConfig.Prod {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic("zap can not init")
	}
	Logger = logger.Sugar()

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), // json格式日志（ELK渲染收集）
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),  // 打印到控制台和文件
		zap.DebugLevel, // 日志级别
	)
	MLogger = zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel), // error级别日志，打印堆栈
	)
}
