package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rashomon/consts"
	"strings"
)

func NoRoute(c *gin.Context) {
	Abort404(c)
}

func Health(c *gin.Context) {
	Success(c, "ok")
}

func Success(c *gin.Context, data interface{}, msg ...string) {
	c.JSON(http.StatusOK, Response{
		Code: consts.Success,
		Msg:  strings.Join(msg, ";"),
		Data: data,
	})
}

// BadRequest 响应 400
func BadRequest(c *gin.Context, err error, msg ...string) {
	c.AbortWithStatusJSON(http.StatusOK, Response{
		Code:  consts.InvalidParams,
		Msg:   strings.Join(msg, ";"),
		Error: err.Error(),
	})
}

// Error 响应
func Error(c *gin.Context, code int, msg ...string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  strings.Join(msg, ";"),
	})
}

func Abort404(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotFound, Response{})
}

func Abort403(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, Response{})
}

func Abort500(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, Response{})
}
