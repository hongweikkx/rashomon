package ctxkv

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"rashomon/consts"
)

type CtxKV struct {
	PlatForm string
}

const CTX_INFO = "info"
const CTX_DEGRADE = "degrade"

func Bind(c *gin.Context) {
	traceId := uuid.New().String()
	ctx := &CtxKV{}
	ctx.PlatForm = c.GetHeader("x-platform")
	c.Set(CTX_INFO, ctx)
	c.Set(consts.TRACE_ID, traceId)
	SetDgd(c, false)
	c.Next()
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
