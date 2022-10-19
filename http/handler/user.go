package handler

import (
	"goserver/global"
	"goserver/service/serviceJwt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type LoginRequest struct {
	Account  string `json:"account" valid:"required~账号不允许为空"`  // 不允许为空
	Password string `json:"password" valid:"required~密码不允许为空"` // 不允许为空
	Type     string `json:"type"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func Login(c *Context, r *LoginRequest) {
	var expiresAt int64 = 0
	if global.Conf.Jwt.ExpireHour > 0 {
		expiresAt = time.Now().Add(time.Duration(global.Conf.Jwt.ExpireHour) * time.Hour).Unix()
	}

	token, err := serviceJwt.CreateToken(serviceJwt.CustomClaims{
		UserId:   "0001",
		UserName: "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	})
	if err != nil {
		c.OutFail(err)
		return
	}

	c.OutSuccess(LoginResponse{
		Token: token,
	})
}

type InfoResponse struct {
	UserName string `json:"user_name"`
}

func Info(c *Context) {
	c.OutSuccess(map[string]interface{}{
		"user_name": c.GetString("user_name"),
	})
}
