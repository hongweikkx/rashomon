package logger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"moul.io/zapgorm2"
	"rashomon/consts"
)

type ZapGorm3 struct {
	zapgorm2.Logger
}

func NewZapGorm3(zapLogger *zap.Logger) ZapGorm3 {
	logger := zapgorm2.New(zapLogger)
	logger.Context = func(ctx context.Context) []zapcore.Field {
		fields := []zap.Field{
			zap.Any(consts.TRACE_ID, ctx.Value(consts.TRACE_ID)),
		}
		return fields
	}
	return ZapGorm3{
		Logger: logger,
	}
}
