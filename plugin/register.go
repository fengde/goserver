package plugin

import (
	"server/plugin/captcha"
	"server/plugin/pprof"
	"server/plugin/prometheus"
)

// 新增的插件，都需要在这里注册登记
func register() {
	store("pprof", pprof.Run)
	store("prometheus", prometheus.Run)
	store("captcha", captcha.Run)
}
