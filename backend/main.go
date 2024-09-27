package main

import (
	"fmt"
	"log"
	"time"

	"daofa/backend/config"
	"daofa/backend/dal"
	"daofa/backend/handler"
	"daofa/backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	config.LoadConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.dbname"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	dal.SetDefault(db)
	handler.InitDB(db) // 初始化 handler 包中的数据库连接

	r := gin.Default()

	// CORS 中间件配置
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true // 动态允许所有源，但会在响应中返回请求的具体源
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "Accept-Language", "Access-Control-Request-Headers", "Access-Control-Request-Method", "Cache-Control", "Connection", "Pragma", "Referer", "Sec-Fetch-Mode", "User-Agent"},
		ExposeHeaders:    []string{"Content-Length", "X-Kuma-Revision"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 提供静态文件服务
	r.Static("/uploads", "./uploads")

	r.POST("/login", handler.Login)

	admin := r.Group("/admin")
	admin.Use(middleware.JWTAuth())
	{
		// 科目相关API
		admin.GET("/subjects", func(c *gin.Context) { handler.GetSubjects(c.Request.Context(), c) })
		admin.POST("/subjects", func(c *gin.Context) { handler.CreateSubject(c.Request.Context(), c) })
		admin.PUT("/subjects/:id", func(c *gin.Context) { handler.UpdateSubject(c.Request.Context(), c) })
		admin.DELETE("/subjects/:id", func(c *gin.Context) { handler.DeleteSubject(c.Request.Context(), c) })

		// 知识点相关API
		admin.GET("/knowledge-points", func(c *gin.Context) { handler.GetKnowledgePoints(c.Request.Context(), c) })
		admin.POST("/knowledge-points", func(c *gin.Context) { handler.CreateKnowledgePoint(c.Request.Context(), c) })
		admin.PUT("/knowledge-points/:id", func(c *gin.Context) { handler.UpdateKnowledgePoint(c.Request.Context(), c) })
		admin.DELETE("/knowledge-points/:id", func(c *gin.Context) { handler.DeleteKnowledgePoint(c.Request.Context(), c) })

		// 习题相关API
		admin.GET("/exercises", func(c *gin.Context) { handler.GetExercises(c.Request.Context(), c) })
		admin.POST("/exercises", func(c *gin.Context) { handler.CreateExercise(c.Request.Context(), c) })
		admin.PUT("/exercises/:id", func(c *gin.Context) { handler.UpdateExercise(c.Request.Context(), c) })
		admin.DELETE("/exercises/:id", func(c *gin.Context) { handler.DeleteExercise(c.Request.Context(), c) })

		// 管理员相关API
		admin.GET("/admins", func(c *gin.Context) { handler.GetAdmins(c.Request.Context(), c) })
		admin.POST("/admins", func(c *gin.Context) { handler.CreateAdmin(c.Request.Context(), c) })
		admin.PUT("/admins/:id", func(c *gin.Context) { handler.UpdateAdmin(c.Request.Context(), c) })
		admin.DELETE("/admins/:id", func(c *gin.Context) { handler.DeleteAdmin(c.Request.Context(), c) })

		// 素材相关API
		admin.GET("/exercise-materials", handler.ListExerciseMaterials)
		admin.POST("/exercise-materials", handler.CreateExerciseMaterial)
		admin.PUT("/exercise-materials/:id", handler.UpdateExerciseMaterial)
		admin.DELETE("/exercise-materials/:id", handler.DeleteExerciseMaterial)
		admin.POST("/upload-image", handler.UploadImage)

		// 新增的题目查询API
		admin.GET("/exercise-questions", handler.GetExerciseQuestions)
	}
	
	// 提供图片加载接口，不需要JWT
	r.StaticFS("/images", gin.Dir("./uploads", true))

	r.Run(fmt.Sprintf(":%d", viper.GetInt("server.port")))
}
