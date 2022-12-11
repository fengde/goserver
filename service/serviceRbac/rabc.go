package serviceRbac

import (
	"context"
	"fmt"
	"goserver/global"

	"github.com/fengde/gocommon/logx"
)

// 检测用户是否拥有权限
func Check(ctx context.Context, userId int64, url string, method string) bool {
	e := *global.Enforcer
	e.LoadPolicy()

	sub := fmt.Sprintf("user_%d", userId)
	obj := url
	act := method

	ok, err := e.Enforce(sub, obj, act)
	if err != nil {
		logx.ErrorWithCtx(ctx, err)
		return false
	}

	if ok {
		return true
	}

	return false
}
