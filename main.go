package main

import (
	"goserver/api"
	"goserver/global"
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

	st := service.NewStartup()
	st.Run()

	safex.Go(api.Run)

	safex.Go(plugin.Run)

	if global.IsDevEnv() {
		safex.Go(test.Run)
	}

	defer safex.Recover(global.Shutdown, st.Shutdown, api.Shutdown, func() {
		logx.Info("bye")
	})

	term := make(chan os.Signal)
	signal.Notify(term, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	logx.Info("signal:", <-term)
}
