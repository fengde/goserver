package serviceSentinel

import (
	"goserver/global"

	sentinel "github.com/alibaba/sentinel-golang/api"
)

func Run() {
	if err := sentinel.InitWithConfigFile(global.Conf.SentinelConfigPath); err != nil {
		panic("sentinel init failed: " + err.Error())
	}
	if err := LoadRules(); err != nil {
		panic("sentinel load rules failed: " + err.Error())
	}
}
