package captcha

import (
	"goserver/http/handler"

	"github.com/fengde/gocommon/captchax"
)

type CaptchaImageResponse struct {
	CaptchaId   string `json:"captcha_id"`
	CaptchaLink string `json:"captcha_link"`
}

func CaptchaImage(c *handler.Context) {
	id, link := captchax.NewCaptchaImage(length)
	c.OutSuccess(CaptchaImageResponse{
		CaptchaId:   id,
		CaptchaLink: link,
	})
}

func CaptchaAudio(c *handler.Context) {
	id, link := captchax.NewCaptchaAudio(length)
	c.OutSuccess(CaptchaImageResponse{
		CaptchaId:   id,
		CaptchaLink: link,
	})
}
