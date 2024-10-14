package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
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

type AlpacaInstruction struct {
	Instruction string `json:"instruction"`
	Input       string `json:"input"`
	Output      string `json:"output"`
}

func main() {
	// 加载配置文件
	loadConfig()

	// 定义命令行参数
	importKnowledgePoints := flag.Bool("import-knowledge-points", false, "导入知识点数据")
	generateInstructions := flag.Bool("generate-instructions", false, "生成 Alpaca 指令集")
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
	} else if *generateInstructions {
		err := generateAlpacaInstructions(1, 3033) // 假设我们处理ID从1到30的习题
		if err != nil {
			log.Fatalf("生成指令集失败: %v", err)
		}
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

		if err := dal.Q.KnowledgePoint.Create(&newPoint); err != nil {
			log.Printf("插入知识点失败: %v", err)
			continue
		}

		if len(point.Children) > 0 {
			insertKnowledgePoints(point.Children, &newPoint.ID)
		}
	}
}

func generateAlpacaInstructions(startID, endID int) error {
	questions, err := dal.Q.Question.Preload(dal.Q.Question.KnowledgePoints).GetQuestionsByIDRange(startID, endID)
	if err != nil {
		return fmt.Errorf("获取习题失败: %v", err)
	}

	instructions := []AlpacaInstruction{}

	for _, q := range questions {
		content := processQuestionContent(q.Content, q.OcrText)
		knowledgePoints := getKnowledgePointsString(q.KnowledgePoints)

		// 1. 习题知识��指令集
		instructionKnowledgePoints := AlpacaInstruction{
			Instruction: "根据以下习题内容，列出相关的知识点。",
			Input:       content,
			Output:      knowledgePoints,
		}
		instructions = append(instructions, instructionKnowledgePoints)

		// 2. 习题内容+知识点做input，解析做output的指令集
		if q.Explanation != nil && len(*q.Explanation) >= 10 {
			instructionExplanation := AlpacaInstruction{
				Instruction: "根据以下习题内容和相关知识点，提供解题思路。",
				Input:       fmt.Sprintf("习题内容：\n%s\n\n相关知识点：\n%s", content, knowledgePoints),
				Output:      *q.Explanation,
			}
			instructions = append(instructions, instructionExplanation)
		}

		// 3. 习题内容+知识点+解析做input，答案做output的指集
		if q.Explanation != nil && len(*q.Explanation) >= 10 && len(q.Answer) >= 10 {
			instructionAnswer := AlpacaInstruction{
				Instruction: "根据以下习题内容、相关知识点和解题思路，给出答案。",
				Input:       fmt.Sprintf("习题内容：\n%s\n\n相关知识点：\n%s\n\n解题思路：\n%s", content, knowledgePoints, *q.Explanation),
				Output:      q.Answer,
			}
			instructions = append(instructions, instructionAnswer)
		}

		// 4. 习题内容做input答案做output的指令集
		if len(q.Answer) >= 10 {
			instructionDirectAnswer := AlpacaInstruction{
				Instruction: "根据以下习题内容，直接给出答案。",
				Input:       content,
				Output:      q.Answer,
			}
			instructions = append(instructions, instructionDirectAnswer)
		}
	}

	// 将指令集转换为JSON并写入文件
	jsonData, err := json.MarshalIndent(instructions, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON编码失败: %v", err)
	}

	err = ioutil.WriteFile("alpaca_instructions.json", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	fmt.Printf("成功生成 %d 条指令\n", len(instructions))
	return nil
}

func getKnowledgePointsString(knowledgePoints []model.KnowledgePoint) string {
	var points []string
	for _, kp := range knowledgePoints {
		points = append(points, kp.Name)
	}
	return strings.Join(points, ", ")
}

func processQuestionContent(content string, ocrText *string) string {
	if ocrText == nil {
		ocrText = new(string)
	}

	// 处理 IMG 标签
	imgRe := regexp.MustCompile(`<img[^>]+src="([^"]+)"[^>]*>`)
	content = imgRe.ReplaceAllStringFunc(content, func(match string) string {
		submatches := imgRe.FindStringSubmatch(match)
		if len(submatches) > 1 {
			if *ocrText != "" {
				return fmt.Sprintf("[图片，内容识别结果：%s]", *ocrText)
			}
			return fmt.Sprintf("[图片：%s]", submatches[1])
		}
		return match
	})

	// 保留允许的标签，移除其他所有 HTML 标签
	allowedTags := []string{"table", "tbody", "tr", "td", "th", "thead", "tfoot", "caption", "colgroup", "col"}
	for _, tag := range allowedTags {
		content = regexp.MustCompile(`(?i)<`+tag+`[^>]*>`).ReplaceAllStringFunc(content, strings.ToLower)
		content = regexp.MustCompile(`(?i)</`+tag+`>`).ReplaceAllStringFunc(content, strings.ToLower)
	}

	// 移除不允许的标签
	content = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(content, "")

	// 将 &nbsp; 实体转换为普通空格
	content = strings.ReplaceAll(content, "&nbsp;", " ")

	// 清理多余的空白字符
	content = strings.TrimSpace(content)
	content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")

	return content
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
