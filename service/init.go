package service

import (
	"goserver/service/serviceSentinel"

	"github.com/fengde/gocommon/logx"
	"github.com/fengde/gocommon/taskx"
)

func Init() {
	go serviceSentinel.Run()
	// //示例
	// go serviceDemo.Run()
}

// 等待子服务正常退出, 每个service维护好自己的Exit，统一执行
func WaitExit() {
	logx.Info("wait services exit...")
	g := taskx.TaskGroup{}
	// // 示例
	// g.Run(serviceDemo.Exit)
	g.Wait()
}
