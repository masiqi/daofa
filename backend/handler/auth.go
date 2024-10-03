package handler

import (
	"net/http"
	"time"

	"daofa/backend/dal"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 1, "msg": err.Error()})
		return
	}

	admin, err := dal.Admin.Where(dal.Admin.Username.Eq(loginData.Username)).First()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": 1, "msg": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": 1, "msg": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin_id": admin.ID,
		"username": admin.Username,
		"exp":      time.Now().Add(viper.GetDuration("JWT_EXPIRATION")).Unix(),
	})

	tokenString, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "", "data": gin.H{"token": tokenString, "username": admin.Username}})
}
