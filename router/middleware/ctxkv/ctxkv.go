package ctxkv

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hongweikkx/rashomon/log"
	"go.uber.org/zap"
)

type CtxKV struct {
	PlatForm string
}

const CTX_INFO = "info"
const CTX_LOG = "logger"
const CTX_DEGRADE = "degrade"

func Bind(c *gin.Context) {
	ctx := &CtxKV{}
	ctx.PlatForm = c.GetHeader("x-platform")
	c.Set(CTX_INFO, ctx)
	c.Set(CTX_LOG, log.MLogger.With(zap.String("trace-id", uuid.New().String())))
	SetDgd(c, false)
	c.Next()
}

func GetInfo(c *gin.Context) *CtxKV {
	if value, isExist := c.Get(CTX_INFO); isExist {
		return value.(*CtxKV)
	}
	return &CtxKV{}
}

func Log(c *gin.Context) *zap.Logger {
	if value, isExist := c.Get(CTX_LOG); isExist {
		return value.(*zap.Logger)
	}
	return log.MLogger
}

func SetDgd(c *gin.Context, dgd bool) {
	c.Set(CTX_DEGRADE, dgd)
}

func GetDgd(c *gin.Context) bool {
	if value, isExist := c.Get(CTX_DEGRADE); isExist {
		return value.(bool)
	}
	return false
}
