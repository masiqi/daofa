package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

func LoadConfig() {
	// 设置配置文件的名称
	viper.SetConfigName(".env")

	// 设置配置文件的类型
	viper.SetConfigType("env")

	// 添加多个可能的配置文件路径
	viper.AddConfigPath(".")
	viper.AddConfigPath("./backend")
	viper.AddConfigPath("../backend")

	// 读取环境变量
	viper.AutomaticEnv()

	// 设置默认值
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 3306)
	viper.SetDefault("DB_USER", "root")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_NAME", "daofa")
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("JWT_EXPIRATION", "720h") // 默认30天
	viper.SetDefault("OCR_URL", "http://10.1.0.242:8000/ocr")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", 6379)
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("未找到配置文件,将使用默认值或环境变量")
		} else {
			panic(fmt.Errorf("读取配置文件时发生错误: %s", err))
		}
	}

	// 打印配置值
	fmt.Println("DB_HOST:", viper.GetString("DB_HOST"))
	fmt.Println("DB_PORT:", viper.GetInt("DB_PORT"))
	fmt.Println("DB_USER:", viper.GetString("DB_USER"))
	fmt.Println("DB_PASSWORD:", viper.GetString("DB_PASSWORD"))
	fmt.Println("DB_NAME:", viper.GetString("DB_NAME"))

	fmt.Println("JWT_SECRET:", viper.GetString("JWT_SECRET"))
	fmt.Println("JWT_EXPIRATION:", viper.GetString("JWT_EXPIRATION"))

	fmt.Println("SERVER_PORT:", viper.GetInt("SERVER_PORT"))

	fmt.Println("OCR_URL:", viper.GetString("OCR_URL"))

	fmt.Println("REDIS_HOST:", viper.GetString("REDIS_HOST"))
	fmt.Println("REDIS_PORT:", viper.GetInt("REDIS_PORT"))
	fmt.Println("REDIS_PASSWORD:", viper.GetString("REDIS_PASSWORD"))
	fmt.Println("REDIS_DB:", viper.GetInt("REDIS_DB"))

	// 打印当前工作目录和配置文件路径(用于调试)
	pwd, _ := os.Getwd()
	fmt.Printf("当前工作目录: %s\n", pwd)
	if viper.ConfigFileUsed() != "" {
		fmt.Printf("使用的配置文件: %s\n", viper.ConfigFileUsed())
	} else {
		fmt.Println("未使用配置文件,使用默认值或环境变量")
	}
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
