package handler

import (
    "github.com/dgrijalva/jwt-go"
    "time"
)

var jwtSecret = []byte("your_jwt_secret") // 确保这里的密钥安全

func generateJWT(username string) (string, error) {
    claims := jwt.MapClaims{}
    claims["username"] = username
    claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // 72小时过期

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}