package dgd

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/router/middleware/ctxkv"
	"github.com/hongweikkx/rashomon/util/enum"
	"go.uber.org/zap"
)

func init() {
	hystrix.DefaultErrorPercentThreshold = 40
	hystrix.DefaultMaxConcurrent = 100
	hystrix.DefaultSleepWindow = 5000
	hystrix.DefaultVolumeThreshold = 25
	hystrix.DefaultTimeout = 5000
}

func Common500() (int, interface{}) {
	return http.StatusServiceUnavailable, ""
}

func dgdName(c *gin.Context) string {
	return c.HandlerName()
}

func Fuse(c *gin.Context, conf *hystrix.CommandConfig, actionF func(ctx context.Context) (int, interface{})) {
	Degrade(c, conf, actionF, Common500)
}

func Degrade(c *gin.Context, conf *hystrix.CommandConfig, actionF func(ctx context.Context) (int, interface{}), degradeF func() (int, interface{})) {
	logger := ctxkv.Log(c)
	if conf != nil {
		hystrix.ConfigureCommand(dgdName(c), *conf)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(hystrix.DefaultTimeout)*time.Millisecond)
	defer cancel()
	doneCh := make(chan bool, 0)
	var status int
	var res interface{}
	go func() {
		err := hystrix.DoC(ctx, dgdName(c), func(ctx context.Context) error {
			defer func() {
				if e := recover(); e != nil {
					logger.Error("[RECOVER FROM PANIC] ",
						zap.Any("error", e),
						zap.String("time", time.Now().Format(time.RFC3339)),
					)
					doneCh <- false
				}
			}()
			status, res = actionF(ctx)
			if status != http.StatusOK {
				return errors.New(enum.ERR_GLOBAL_COMMON)
			}
			return nil
		}, nil)
		doneCh <- err == nil
	}()
	done := <-doneCh
	if !done {
		ctxkv.SetDgd(c, true)
		status, res = degradeF()
	}
	if ress, ok := res.(string); ok {
		c.String(status, ress)
	} else {
		c.JSON(status, res)
	}
}
