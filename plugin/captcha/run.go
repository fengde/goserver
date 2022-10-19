package captcha

import (
	"goserver/http"
	"goserver/http/handler"

	"github.com/fengde/gocommon/captchax"
	"github.com/gin-gonic/gin"
)

var length int

func Run(setting *Setting) {
	length = setting.Length
	if length <= 0 {
		length = 4
	}
	engine := http.GetGinEngine()
	// 生成图形验证码
	engine.POST("/api/captcha/image", handler.WrapF(CaptchaImage))
	// 生成音频验证码
	engine.POST("/api/captcha/audio", handler.WrapF(CaptchaAudio))
	// 访问验证码资源
	engine.GET("/captcha/*path", gin.WrapH(captchax.LinkHandle()))
}
