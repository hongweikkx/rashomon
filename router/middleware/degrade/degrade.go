package dgd

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/log"
	"github.com/hongweikkx/rashomon/util/enum"
)

const (
	DEGRADE = "degrade"
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
					log.SugarLogger.Error("[RECOVER] ", e, " [STACK] ", string(debug.Stack()))
					doneCh <- false
				}
			}()
			status, res = actionF(ctx)
			if status != http.StatusOK {
				log.SugarLogger.Errorf("%s: action:%s, status:%d, res:%+v", enum.ERR_GLOBAL_THIRD_PARTY,
					runtime.FuncForPC(reflect.ValueOf(actionF).Pointer()).Name(), status, res)
				return errors.New(enum.ERR_GLOBAL_THIRD_PARTY)
			}
			return nil
		}, nil)
		doneCh <- err == nil
	}()
	done := <-doneCh
	if !done {
		status, res = degradeF()
	}
	if ress, ok := res.(string); ok {
		c.String(status, ress)
	} else {
		c.JSON(status, res)
	}
}
