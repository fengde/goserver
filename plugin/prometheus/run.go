package prometheus

import (
	"net/http"

	_prometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var reg = _prometheus.NewRegistry()

// 启动
func Run(setting *Setting) {

	reg.MustRegister(NewSystemCollectors()...)
	reg.MustRegister(NewHttpCollectors()...)

	http.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))

	http.ListenAndServe(setting.Address, nil)
}
