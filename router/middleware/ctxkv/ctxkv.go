package ctxkv

import (
	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/log"
)

type CtxKV struct {
	PlatForm string
}

const CTX_INFO = "info"

func Bind(c *gin.Context) {
	ctx := &CtxKV{}
	ctx.PlatForm = c.GetHeader("x-platform")
	c.Set(CTX_INFO, ctx)
	c.Next()
}

func Get(c *gin.Context) *CtxKV {
	value, isExist := c.Get(CTX_INFO)
	if !isExist {
		log.SugarLogger.Errorf("get ctx info not valid")
		return &CtxKV{}
	}
	return value.(*CtxKV)
}
