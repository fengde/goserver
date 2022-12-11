package serviceJwt

import (
	"goserver/global"

	"github.com/dgrijalva/jwt-go"
	"github.com/fengde/gocommon/errorx"
)

type CustomClaims struct {
	UserId int64 `json:"user_id"`
	jwt.StandardClaims
}

// 解析token
func ParseJwt(token string) (*CustomClaims, error) {
	obj, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Conf.Jwt.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if t, ok := obj.Claims.(*CustomClaims); ok && obj.Valid {
		return t, nil
	}

	return nil, errorx.New("token is valid")
}

// 创建token
func CreateToken(claims CustomClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(global.Conf.Jwt.Secret))
}
