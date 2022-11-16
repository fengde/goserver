package global

import (
	"context"
	"goserver/consts"
	"time"
)

// 当前运行环境是否为dev
func IsDevEnv() bool {
	return Conf.Env == consts.ENV_DEV
}

// 当前运行环境是否为qa
func IsQaEnv() bool {
	return Conf.Env == consts.ENV_QA
}

// 当前运行环境是否为online
func IsOnlineEnv() bool {
	return Conf.Env == consts.ENV_ONLINE
}

// 系统是否还在运行，带sleep
func Continue(sleep ...time.Duration) bool {
	if len(sleep) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), sleep[0])
		defer cancel()

		select {
		case <-exit:
			return false
		case <-ctx.Done():
			return true
		}
	}

	select {
	case <-exit:
		return false
	default:
		return true
	}
}

// 设置系统关闭
func Shutdown() {
	close(exit)
}
