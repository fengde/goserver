package serviceDemo

import (
	"server/global"
	"time"

	"github.com/fengde/gocommon/logx"
	"github.com/fengde/gocommon/safex"
)

var exit = make(chan bool, 1)

func Run() {

	defer safex.Recover(func() { exit <- true })

	for global.Continue() {
		safex.Func(func() {
			logx.Info("serviceDemo Run ......")
		})

		if !global.Continue() {
			break
		}

		time.Sleep(time.Second)
	}
}

func Exit() {
	<-exit
}
