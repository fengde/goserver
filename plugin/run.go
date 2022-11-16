package plugin

import (
	"fmt"
	"goserver/global"
	"goserver/plugin/captcha"
	"goserver/plugin/pprof"
	"goserver/plugin/prometheus"
	"reflect"

	"github.com/fengde/gocommon/jsonx"
	"github.com/fengde/gocommon/logx"
	"github.com/fengde/gocommon/safex"
)

var registers = map[string]any{}

// 注册插件
func register(name string, handler any) {
	if _, ok := registers[name]; !ok {
		panic("plugin register again!")
	}
	registers[name] = handler
}

// 启动运行, 新增的插件，都需要在这里注册登记
func Run() {
	register("pprof", pprof.Run)
	register("prometheus", prometheus.Run)
	register("captcha", captcha.Run)

	start()
}

func start() {
	for i := range global.Conf.Plugins {
		func(plg global.Plugin) {
			if plg.Open {
				if handler, ok := registers[plg.Name]; ok {
					handlerType := reflect.TypeOf(reflect.ValueOf(handler).Interface())
					argNum := handlerType.NumIn()
					args := make([]reflect.Value, argNum)
					for i := 0; i < argNum; i++ {
						paramPtr := handlerType.In(i)
						paramKind := paramPtr.Kind()
						if paramKind == reflect.Ptr {
							argi := reflect.New(paramPtr.Elem()).Interface()
							if err := jsonx.UnmarshalString(plg.Setting, argi); err != nil {
								panic(fmt.Sprintf("plugin: [%s] the setting of plugin unmarshal error: %s", plg.Name, err.Error()))
							}
							args[i] = reflect.ValueOf(argi)
						}
					}
					safex.Go(func() {
						logx.Info("start plugin: ", plg.Name)
						reflect.ValueOf(handler).Call(args)
					})
				}
			}
		}(global.Conf.Plugins[i])
	}
}
