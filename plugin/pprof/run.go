package pprof

import (
	"server/http"
	"strings"

	"github.com/fengde/gocommon/filex"
	"github.com/fengde/gocommon/logx"
	"github.com/fengde/gocommon/toolx"
	"github.com/gin-contrib/pprof"
)

func Run(setting *Setting) {
	pprofUrl := "/pprof/" + toolx.NewCharCode(32)

	if err := filex.Write(setting.RuntimeDataPath, strings.NewReader("pprof_url: "+pprofUrl)); err != nil {
		logx.Error("pprof filex.Write", err)
		return
	}

	logx.Info("pprof path: ", pprofUrl)

	pprof.Register(http.GetGinEngine(), pprofUrl)
}
