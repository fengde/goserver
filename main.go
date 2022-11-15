package main

import (
	"goserver/global"
	"goserver/http"
	"goserver/plugin"
	"goserver/service"
	"goserver/test"
	"os"
	"os/signal"
	"syscall"

	"github.com/fengde/gocommon/logx"
	"github.com/fengde/gocommon/safex"
)

func main() {
	if err := global.Init(); err != nil {
		panic(err)
	}

	svc := service.NewService()
	svc.Run()

	defer safex.Recover(global.Exist, svc.Close, func() {
		logx.Info("bye bye")
	})

	safex.Go(http.Run)

	safex.Go(plugin.Run)

	safex.Go(test.Start)

	listenSignal()

	if err := http.Shutdown(); err != nil {
		logx.Error(err)
	}
}

func listenSignal() {
	term := make(chan os.Signal)
	signal.Notify(term, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	logx.Info("get signal:", <-term)
}
