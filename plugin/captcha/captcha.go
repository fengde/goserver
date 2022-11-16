package captcha

import (
	"goserver/api/handler"

	"github.com/fengde/gocommon/captchax"
)

type CaptchaImageResponse struct {
	CaptchaId   string `json:"captcha_id"`
	CaptchaLink string `json:"captcha_link"`
}

func CaptchaImage(c *handler.Context) (*CaptchaImageResponse, error) {
	id, link := captchax.NewCaptchaImage(length)
	return &CaptchaImageResponse{
		CaptchaId:   id,
		CaptchaLink: link,
	}, nil
}

func CaptchaAudio(c *handler.Context) (*CaptchaImageResponse, error) {
	id, link := captchax.NewCaptchaAudio(length)
	return &CaptchaImageResponse{
		CaptchaId:   id,
		CaptchaLink: link,
	}, nil
}
