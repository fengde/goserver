package serviceSentinel

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
)

// 调用封装
func Entry(resource string, normalHandler func(), blockHandler func(b *base.BlockError), opts ...sentinel.EntryOption) {
	e, b := sentinel.Entry(resource, opts...)
	if b != nil {
		blockHandler(b)
	} else {
		defer e.Exit()
		normalHandler()
	}
}
