package service

import (
	"goserver/service/serviceSentinel"

	"github.com/fengde/gocommon/taskx"
)

var tg = taskx.NewTaskGroup()

// 这里开机启动运行服务
func Run() {
	tg.Run(serviceSentinel.Run)
}

// 统一等待结束
func Shutdown() {
	tg.Wait()
}
