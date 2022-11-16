package service

import (
	"goserver/global"
	"time"

	"github.com/fengde/gocommon/safex"
)

// 带函数锁执行函数，用于分布式场景，避免多开情况下的并发冲突
func LoopWrapfWithLocker(fn func(), interval time.Duration, lockId string, lockTimeout time.Duration) func() {
	if interval == 0 || lockId == "" || lockTimeout/time.Second < 1 {
		panic("入参非法")
	}
	return func() {
		for {
			global.Locker.Lock(lockId, int64(lockTimeout/time.Second), fn)
			if !global.Continue(interval) {
				break
			}
		}
	}
}

// 循环执行函数
func LoopWrapf(fn func(), interval time.Duration) func() {
	if interval == 0 {
		panic("入参非法")
	}
	return func() {
		for {
			safex.Func(fn)
			if !global.Continue(interval) {
				break
			}
		}
	}
}
