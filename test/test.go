package test

import (
	"goserver/global"
	"goserver/service/serviceUser"

	"github.com/fengde/gocommon/jsonx"
	"github.com/fengde/gocommon/logx"
)

// 内部逻辑测试
func Run() {
	// test1209_GORM()
}

func test1208_RBAC() {
	global.Enforcer.AddPolicy("admin", "/api/user/info", "POST")
	for _, row := range global.Enforcer.GetPolicy() {
		logx.Info(row)
	}
}

func test1209_GORM() {
	u := serviceUser.GetUserInfo(1)
	logx.Info(jsonx.MarshalToStringNoErr(u))

	roles, err := serviceUser.GetUserRoles(1)
	logx.Info(jsonx.MarshalToStringNoErr(roles), err)
}
