package test

import "server/global"

// 内部逻辑测试，仅仅在dev环境打开
func Start() {
	if !global.IsDevEnv() {
		return
	}
	hereStartYourTest()
}
