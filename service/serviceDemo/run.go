package serviceDemo

import (
	"goserver/global"
	"time"
)

func Run() {
	for {
		// todo something
		if !global.Continue(time.Second) {
			break
		}
	}
}
