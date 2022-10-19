package http

import (
	"context"
	"log"
	"net/http"
	"server/global"
	"server/http/router"

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

func Shutdown() error {
	return srv.Shutdown(context.Background())
}

func GetGinEngine() *gin.Engine {
	return engine
}
