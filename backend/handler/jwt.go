package handler

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

func generateJWT(username string) (string, error) {
	// 从配置中读取 JWT 密钥
	jwtSecret := []byte(viper.GetString("JWT_SECRET"))

	// 从配置中读取 JWT 过期时间
	jwtExpiration := viper.GetDuration("JWT_EXPIRATION")

	claims := jwt.MapClaims{}
	claims["username"] = username
	claims["exp"] = time.Now().Add(jwtExpiration).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
