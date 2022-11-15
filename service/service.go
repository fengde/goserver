package service

import (
	"goserver/service/serviceSentinel"

	"github.com/fengde/gocommon/taskx"
)

type Service struct {
	g *taskx.TaskGroup
}

func NewService() *Service {
	return &Service{
		g: taskx.NewTaskGroup(),
	}
}

// 这里运行服务
func (p *Service) Run() {
	p.g.Run(serviceSentinel.Run)
}

// 统一等待关闭
func (p *Service) Close() {
	p.g.Wait()
}
