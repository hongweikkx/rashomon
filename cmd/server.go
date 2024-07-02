package cmd

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"rashomon/conf"
	"rashomon/consts"
	"rashomon/pkg/logger"
	"rashomon/router"
	"syscall"
	"time"
)

var WebServer = &cobra.Command{
	Use:   "server",
	Short: "Start Web Server",
	Run:   runWebServer,
	Args:  cobra.NoArgs,
}

func runWebServer(cmd *cobra.Command, args []string) {
	if conf.AppConfig.Service.AppMode == consts.APP_MOD_PROD {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.New()
	router.Router(engine)
	serv := &http.Server{
		Addr:           conf.AppConfig.Service.HttpPort,
		Handler:        engine,
		ReadTimeout:    time.Duration(conf.AppConfig.Http.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(conf.AppConfig.Http.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << conf.AppConfig.Http.MaxHeaderBytes,
	}

	go func() {
		if err := serv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) && err != nil {
			logger.Fatal(context.Background(), "http serv err:", zap.Error(err))
		}
	}()
	logger.Info(context.Background(), "server started on ...", zap.Any("addr", conf.AppConfig.Service.HttpPort))
	// wait to stop
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit
	serv.SetKeepAlivesEnabled(false)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := serv.Shutdown(ctx); err != nil {
		logger.Error(ctx, "http server forced to shutdown:", zap.Error(err))
	}
	logger.Info(context.Background(), "server exit.")
}

// 加载全局中间件
//func registerGlobalMiddleWare(router *gin.Engine) {
//	router.Use(
//		common.Recovery(),
//		common._Logger(),
//	)
//	store := cookie.NewStore([]byte(conf.CookieSecret))
//	router.Use(common.Cors())
//	router.Use(sessions.Sessions(conf.SessionName, store))
//}
