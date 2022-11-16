package startup

import (
	"goserver/global"
	"goserver/service/serviceSentinel"
	"time"

	"github.com/fengde/gocommon/safex"
	"github.com/fengde/gocommon/taskx"
)

type Startup struct {
	g *taskx.TaskGroup
}

func NewStartup() *Startup {
	return &Startup{
		g: taskx.NewTaskGroup(),
	}
}

// 这里开机启动运行服务
func (p *Startup) Run() {
	p.g.Run(serviceSentinel.Run)
}

// 带函数锁执行函数，用于分布式场景，避免多开情况下，并发冲突
func (p *Startup) LoopWrapfWithLocker(fn func(), sleep time.Duration, lockId string, lockTimeout time.Duration) func() {
	return func() {
		for {
			global.Locker.Lock(lockId, int64(lockTimeout/time.Second), fn)
			if !global.Continue(sleep) {
				break
			}
		}
	}
}

// 循环执行函数
func (p *Startup) LoopWrapf(fn func(), sleep time.Duration) func() {
	return func() {
		for {
			safex.Func(fn)
			if !global.Continue(sleep) {
				break
			}
		}
	}
}

// 统一等待关闭
func (p *Startup) Close() {
	p.g.Wait()
}
