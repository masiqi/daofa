package main

import (
    "daofa/backend/dal"
    "daofa/backend/model"
    "fmt"
    "log"
    "os"

    "github.com/spf13/viper"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: go run cmd/create_admin.go <username> <password>")
        os.Exit(1)
    }

    username := os.Args[1]
    password := os.Args[2]

    // 加载配置
    viper.SetConfigFile("config.yaml")
    if err := viper.ReadInConfig(); err != nil {
        log.Fatalf("Error reading config file: %s", err)
    }

    // 连接数据库
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        viper.GetString("database.username"),
        viper.GetString("database.password"),
        viper.GetString("database.host"),
        viper.GetInt("database.port"),
        viper.GetString("database.dbname"))

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    dal.SetDefault(db)

    // 创建管理员
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Fatalf("Failed to hash password: %v", err)
    }

    admin := model.Admin{
        Username: username,
        Password: string(hashedPassword),
    }

    if err := dal.Q.Admin.Create(&admin); err != nil {
        log.Fatalf("Failed to create admin: %v", err)
    }

    fmt.Printf("Admin created successfully: %s\n", username)
}