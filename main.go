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
	"time"

	"github.com/fengde/gocommon/logx"
	"github.com/fengde/gocommon/safex"
)

func main() {
	if err := global.Init(); err != nil {
		panic(err)
	}

	service.Init()

	defer safex.Recover(global.Exist, service.WaitExit, func() {
		logx.Info("bye bye")
	})

	safex.Go(http.Run)

	safex.Go(plugin.Run)

	safex.Go(test.Start)

	listenSignal()

	if err := http.Shutdown(); err != nil {
		logx.Error(err)
	}

	time.Sleep(time.Second)
}

func listenSignal() {
	term := make(chan os.Signal)
	signal.Notify(term, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	logx.Info("get signal:", <-term)
}
