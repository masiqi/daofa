package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"daofa/backend/dal"
	"daofa/backend/handler"
	"daofa/backend/middleware"
	"daofa/backend/queue"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"daofa/backend/model"
)

var redisClient *redis.Client

type KnowledgePointJSON struct {
	ID       int                  `json:"id"`
	Name     string               `json:"name"`
	IsLeaf   bool                 `json:"is_leaf"`
	Level    int                  `json:"level"`
	Children []KnowledgePointJSON `json:"children,omitempty"`
}

func main() {
	// 加载配置文件
	loadConfig()

	// 定义命令行参数
	importKnowledgePoints := flag.Bool("import-knowledge-points", false, "导入知识点数据")
	flag.Parse()

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_HOST"),
		viper.GetInt("DB_PORT"),
		viper.GetString("DB_NAME"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 初始化 DAL
	dal.SetDefault(db)

	if *importKnowledgePoints {
		importKnowledgePointsFromJSON()
	} else {
		// 其他主程序逻辑
		runServer()
	}
}

func loadConfig() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// 设置默认值
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("DB_PORT", 3306)
	viper.SetDefault("REDIS_PORT", 6379)
	viper.SetDefault("REDIS_DB", 0)
}

func importKnowledgePointsFromJSON() {
	// 读取JSON文件
	jsonFile, err := os.Open("sql/knowledge_points.json")
	if err != nil {
		log.Fatalf("无法打开JSON文件: %v", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var knowledgePoints []KnowledgePointJSON
	json.Unmarshal(byteValue, &knowledgePoints)

	// 插入知识点
	insertKnowledgePoints(knowledgePoints, nil)

	fmt.Println("知识点数据导入完成")
}

func insertKnowledgePoints(points []KnowledgePointJSON, parentID *int32) {
	for _, point := range points {
		newPoint := model.KnowledgePoint{
			Name:      point.Name,
			IsLeaf:    point.IsLeaf,
			ParentID:  parentID,
			SubjectID: 1, // 假设所有知识点属于同一个学科,实际使用时可能需要调整
		}

		err := dal.Q.KnowledgePoint.CreateKnowledgePoint(newPoint.SubjectID, newPoint.ParentID, newPoint.Name, nil, newPoint.IsLeaf)
		if err != nil {
			log.Printf("插入知识点失败: %v", err)
			continue
		}

		// 获取刚插入的知识点ID
		insertedPoint, err := dal.Q.KnowledgePoint.GetKnowledgePointByName(newPoint.Name)
		if err != nil {
			log.Printf("获取插入的知识点失败: %v", err)
			continue
		}

		newID := insertedPoint.ID

		// 递归插入子节点
		if len(point.Children) > 0 {
			insertKnowledgePoints(point.Children, &newID)
		}
	}
}

func runServer() {
	// Redis 连接
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("REDIS_HOST"), viper.GetInt("REDIS_PORT")),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       viper.GetInt("REDIS_DB"),
	})

	// 测试 Redis 连接
	_, err := redisClient.Ping(redisClient.Context()).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	// 初始化 Redis
	queue.InitRedis(redisClient)

	// 启动习题处理器
	go handler.ProcessPendingQuestions(redisClient.Context())

	// 启动图片OCR任务处理器
	go handler.ProcessImageOCRTasks(redisClient.Context())

	r := gin.Default()

	// CORS 中间件配置
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true // 动态允许所有源，但会在响应中返回请求的具体
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
		admin.GET("/subjects", handler.GetSubjects)
		admin.POST("/subjects", handler.CreateSubject)
		admin.PUT("/subjects/:id", handler.UpdateSubject)
		admin.DELETE("/subjects/:id", handler.DeleteSubject)

		// 知识点相关API
		admin.GET("/knowledge-points", handler.GetKnowledgePoints)
		admin.POST("/knowledge-points", handler.CreateKnowledgePoint)
		admin.PUT("/knowledge-points/:id", handler.UpdateKnowledgePoint)
		admin.DELETE("/knowledge-points/:id", handler.DeleteKnowledgePoint)
		admin.GET("/knowledge-points/:id", handler.GetKnowledgePoint) // 新增这一行

		// 管理员相关API
		admin.GET("/admins", handler.GetAdmins)
		admin.POST("/admins", handler.CreateAdmin)
		admin.PUT("/admins/:id", handler.UpdateAdmin)
		admin.DELETE("/admins/:id", handler.DeleteAdmin)

		// 新增的路由
		admin.POST("/enqueue-questions", handler.EnqueueQuestions)
		admin.GET("/queue-status", handler.GetQueueStatus)

		// 新增题目相关的路由
		admin.POST("/questions", handler.CreateQuestion)
		admin.GET("/questions", handler.ListQuestions)
		admin.GET("/questions/:id", handler.GetQuestion)
		admin.PUT("/questions/:id", handler.UpdateQuestion)
		admin.DELETE("/questions/:id", handler.DeleteQuestion)
		admin.GET("/questions/search", handler.SearchQuestions)

		// 目知识点关联的路由
		admin.POST("/questions/:id/knowledge-points", handler.AddQuestionKnowledgePoint)
		admin.DELETE("/questions/:id/knowledge-points/:knowledge_point_id", handler.RemoveQuestionKnowledgePoint)

		// 题目类型关由
		admin.POST("/question-types", handler.CreateQuestionType)
		admin.GET("/question-types", handler.ListQuestionTypes)
		admin.GET("/question-types/:id", handler.GetQuestionType)
		admin.PUT("/question-types/:id", handler.UpdateQuestionType)
		admin.DELETE("/question-types/:id", handler.DeleteQuestionType)

		// 添加 image_ocr_tasks 相关的路由
		admin.GET("/image-ocr-tasks", handler.ListImageOCRTasks)
		admin.GET("/image-ocr-tasks/:id", handler.GetImageOCRTask)
		admin.POST("/image-ocr-tasks", handler.CreateImageOCRTask)
		admin.PUT("/image-ocr-tasks/:id", handler.UpdateImageOCRTask)
		admin.DELETE("/image-ocr-tasks/:id", handler.DeleteImageOCRTask)
		admin.GET("/image-ocr-tasks/search", handler.SearchImageOCRTasks)
	}

	// 提供图片加载接口，不需要JWT
	r.StaticFS("/images", gin.Dir("./uploads", true))

	// 添加新的路由,不需要JWT认证
	r.POST("/enqueue-image-ocr", handler.EnqueueImageOCR)

	r.Run(fmt.Sprintf(":%d", viper.GetInt("SERVER_PORT")))
}