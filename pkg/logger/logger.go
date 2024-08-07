package logger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"rashomon/conf"
	"rashomon/consts"
)

var _Logger *zap.Logger

// 负责设置 encoding 的日志格式
func getEncoder() zapcore.Encoder {
	// 获取一个指定的的EncoderConfig，进行自定义
	encodeConfig := zap.NewProductionEncoderConfig()
	// 序列化时间。eg: 2022-09-01T19:11:35.921+0800
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// "time":"2022-09-01T19:11:35.921+0800"
	encodeConfig.TimeKey = "time"
	// 将Level序列化为全大写字符串。例如，将info level序列化为INFO。
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encodeConfig.LineEnding = zapcore.DefaultLineEnding
	return zapcore.NewJSONEncoder(encodeConfig)
}

// 负责日志写入的位置
func getLogWriter(filename string, maxsize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  // 文件位置
		MaxSize:    maxsize,   // 进行切割之前,日志文件的最大大小(MB为单位)
		MaxAge:     maxAge,    // 保留旧文件的最大天数
		MaxBackups: maxBackup, // 保留旧文件的最大个数
		Compress:   false,     // 是否压缩/归档旧文件
	}
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
}

func Init() {
	// 获取日志写入位置
	writeSyncer := getLogWriter(conf.AppConfig.Log.FileName, conf.AppConfig.Log.MaxSize, conf.AppConfig.Log.MaxBackups, conf.AppConfig.Log.MaxAge)
	// 获取日志编码格式
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zap.DebugLevel)
	_Logger = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(_Logger)
	return
}

func Quit() {
	_Logger.Sync()
}

func Debug(c context.Context, moduleName string, fields ...zap.Field) {
	fields = append(fields, zap.Any(consts.TRACE_ID, c.Value(consts.TRACE_ID)))
	_Logger.Debug(moduleName, fields...)
}

// Info 告知类日志
func Info(c context.Context, moduleName string, fields ...zap.Field) {
	fields = append(fields, zap.Any(consts.TRACE_ID, c.Value(consts.TRACE_ID)))
	_Logger.Info(moduleName, fields...)
}

// Warn 警告类
func Warn(c context.Context, moduleName string, fields ...zap.Field) {
	fields = append(fields, zap.Any(consts.TRACE_ID, c.Value(consts.TRACE_ID)))
	_Logger.Warn(moduleName, fields...)
}

// Error 错误时记录，不应该中断程序，查看日志时重点关注
func Error(c context.Context, moduleName string, fields ...zap.Field) {
	fields = append(fields, zap.Any(consts.TRACE_ID, c.Value(consts.TRACE_ID)))
	_Logger.Error(moduleName, fields...)
}

// Fatal 级别同 Error(), 写完 log 后调用 os.Exit(1) 退出程序
func Fatal(c context.Context, moduleName string, fields ...zap.Field) {
	fields = append(fields, zap.Any(consts.TRACE_ID, c.Value(consts.TRACE_ID)))
	_Logger.Fatal(moduleName, fields...)
}
