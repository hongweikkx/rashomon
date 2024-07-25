package ctxkv

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"rashomon/consts"
)

type CtxKV struct {
	PlatForm string
}

func Bind(c *gin.Context) {
	traceId := uuid.New().String()
	ctx := &CtxKV{}
	ctx.PlatForm = c.GetHeader("x-platform")
	c.Set(consts.CTX_INFO, ctx)
	c.Set(consts.TRACE_ID, traceId)
	c.Next()
}

func GetCtxInfo(c *gin.Context) (*CtxKV, bool) {
	ctxInfo, ok := c.Get(consts.CTX_INFO)
	if !ok {
		return nil, false
	}
	return ctxInfo.(*CtxKV), true
}
