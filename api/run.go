package api

import (
	"context"
	"goserver/api/router"
	"goserver/global"
	"log"
	"net/http"

	"github.com/fengde/gocommon/logx"
	"github.com/gin-gonic/gin"
)

var srv *http.Server
var engine = gin.Default()

func Run() {
	router.Init(engine)

	logx.Info("listen on", global.Conf.HttpAddress)

	srv = &http.Server{
		Addr:    global.Conf.HttpAddress,
		Handler: engine,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func Shutdown() {
	srv.Shutdown(context.Background())
}

func GetGinEngine() *gin.Engine {
	return engine
}
