package dashboard

import (
	"context"
	"github.com/hongweikkx/rashomon/dashboard/route"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/log"
)

// Dashboard
type Dashboard struct {
	Serv *http.Server
}

var DashboardIns *Dashboard

func Start() error {
	if conf.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	err := route.Router(engine)
	if err != nil {
		return err
	}
	serv := &http.Server{
		Addr:    conf.AppConfig.DashBoard.Addr,
		Handler: engine,
	}
	DashboardIns = &Dashboard{Serv: serv}
	go func() {
		if err := serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.SugarLogger.Fatal("http serv err:", err.Error())
		}
	}()
	return nil
}

func Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := DashboardIns.Serv.Shutdown(ctx); err != nil {
		log.SugarLogger.Info("http server forced to shutdown:", err)
	}
}

