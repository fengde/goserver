package service

import (
	"goserver/global"
	"time"

	"github.com/fengde/gocommon/safex"
)

// 带函数锁执行函数，用于分布式场景，避免多开情况下，并发冲突
func LoopWrapfWithLocker(fn func(), sleep time.Duration, lockId string, lockTimeout time.Duration) func() {
	return func() {
		for {
			global.Locker.Lock(lockId, int64(lockTimeout/time.Second), fn)
			if !global.Continue(sleep) {
				break
			}
		}
	}
}

// 循环执行函数
func LoopWrapf(fn func(), sleep time.Duration) func() {
	return func() {
		for {
			safex.Func(fn)
			if !global.Continue(sleep) {
				break
			}
		}
	}
}
