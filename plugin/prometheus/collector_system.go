package prometheus

import (
	"goserver/global"

	"github.com/fengde/gocommon/timex"
	_prometheus "github.com/prometheus/client_golang/prometheus"
)

// 系统相关参数
func NewSystemCollectors() []_prometheus.Collector {
	counter := _prometheus.NewCounterFunc(_prometheus.CounterOpts{
		Namespace: "",
		Subsystem: "",
		Name:      "uptime",
		Help:      "运行时长（秒）",
	}, func() float64 {
		return float64(timex.NowUnix() - global.StartUnix)
	})

	return []_prometheus.Collector{counter}
}
