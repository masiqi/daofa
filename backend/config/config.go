package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func LoadConfig() {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// 设置配置
	viper.SetDefault("server.port", getEnvAsInt("SERVER_PORT", 8080))

	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.host", os.Getenv("DB_HOST"))
	viper.SetDefault("database.port", getEnvAsInt("DB_PORT", 3306))
	viper.SetDefault("database.username", os.Getenv("DB_USER"))
	viper.SetDefault("database.password", os.Getenv("DB_PASSWORD"))
	viper.SetDefault("database.dbname", os.Getenv("DB_NAME"))

	viper.SetDefault("jwt.secret", os.Getenv("JWT_SECRET"))
	viper.SetDefault("jwt.expiration", os.Getenv("JWT_EXPIRATION"))
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}