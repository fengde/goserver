package startup

import (
	"goserver/service/serviceSentinel"

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

// 统一等待关闭
func (p *Startup) Close() {
	p.g.Wait()
}
