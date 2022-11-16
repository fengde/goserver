package captcha

import (
	http "goserver/api"
	"goserver/api/router"

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
	// 访问验证码资源
	engine.GET("/captcha/*path", gin.WrapH(captchax.LinkHandle()))
	// 生成图形验证码
	router.POST(&engine.RouterGroup, "/api/captcha/image", CaptchaImage)
	// 生成音频验证码
	router.POST(&engine.RouterGroup, "/api/captcha/audio", CaptchaAudio)

}
