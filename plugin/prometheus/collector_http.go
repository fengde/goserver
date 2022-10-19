package prometheus

import (
	"server/service/serviceHttp"

	_prometheus "github.com/prometheus/client_golang/prometheus"
)

// Http相关信息采集
func NewHttpCollectors() []_prometheus.Collector {
	counter := _prometheus.NewCounterFunc(_prometheus.CounterOpts{
		Namespace: "",
		Subsystem: "",
		Name:      "http_request_total",
		Help:      "总http请求数",
	}, func() float64 {
		info := serviceHttp.GetCountInfo()
		return float64(info.Total)
	})

	counter2 := _prometheus.NewGaugeFunc(_prometheus.GaugeOpts{
		Namespace: "",
		Subsystem: "",
		Name:      "http_request_current",
		Help:      "当前http请求数",
	}, func() float64 {
		info := serviceHttp.GetCountInfo()
		return float64(info.Current)
	})

	counter3 := _prometheus.NewCounterFunc(_prometheus.CounterOpts{
		Namespace: "",
		Subsystem: "",
		Name:      "http_request_total_200",
		Help:      "总的正常请求数",
	}, func() float64 {
		info := serviceHttp.GetCountInfo()
		return float64(info.Total200)
	})

	counter4 := _prometheus.NewCounterFunc(_prometheus.CounterOpts{
		Namespace: "",
		Subsystem: "",
		Name:      "http_request_size_bytes",
		Help:      "总请求字节数",
	}, func() float64 {
		info := serviceHttp.GetCountInfo()
		return float64(info.RequestSizeBytes)
	})

	counter5 := _prometheus.NewCounterFunc(_prometheus.CounterOpts{
		Namespace: "",
		Subsystem: "",
		Name:      "http_response_size_bytes",
		Help:      "总请求返回字节数",
	}, func() float64 {
		info := serviceHttp.GetCountInfo()
		return float64(info.ResponseSizeBytes)
	})

	return []_prometheus.Collector{counter, counter2, counter3, counter4, counter5}
}
