package plugin

import (
	"fmt"
	"goserver/conf"
	"goserver/global"
	"reflect"

	"github.com/fengde/gocommon/jsonx"
	"github.com/fengde/gocommon/logx"
	"github.com/fengde/gocommon/safex"
)

var registers = map[string]any{}

// 注册插件
func store(name string, handler any) {
	if _, ok := registers[name]; ok {
		panic("plugin register again!")
	}
	registers[name] = handler
}

// 启动运行
func Run() {
	register()
	start()
}

func start() {
	for i := range global.Conf.Plugins {
		func(plg conf.Plugin) {
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
