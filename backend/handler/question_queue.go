package handler

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"daofa/backend/dal"
	"daofa/backend/model"
	"daofa/backend/queue"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

// EnqueueQuestions 将题目加入队列
func EnqueueQuestions(c *gin.Context) {
	// 获取科目ID
	subjectIDStr := c.PostForm("subjectId")
	subjectID, err := strconv.Atoi(subjectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的科目ID"})
		return
	}

	// 检查科目是否存在
	subject, err := dal.Q.Subject.GetSubjectByID(int32(subjectID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "科目不存在"})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法获取上传的文件"})
		return
	}

	// 读取文件内容
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法打开文件"})
		return
	}
	defer fileContent.Close()

	jsonData, err := ioutil.ReadAll(fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取文件内容"})
		return
	}

	// 解析JSON数据
	var questions []queue.QuestionItem
	err = json.Unmarshal(jsonData, &questions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的JSON格式"})
		return
	}

	// 将科目ID添加到每个问题中
	for i := range questions {
		questions[i].SubjectID = int32(subjectID)
	}

	// 将问题加入队列
	ctx := context.Background()
	for _, q := range questions {
		err = queue.EnqueueQuestion(ctx, q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "将题目加入队列失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%d 个题目已成功加入队列", len(questions)), "subject": subject})
}

// GetQueueStatus 获取队列状态
func GetQueueStatus(c *gin.Context) {
	ctx := context.Background()
	questionQueueStatus, err := queue.GetQueueStatus(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取队列状态失败"})
		return
	}

	imageOCRQueueLength, err := queue.GetImageOCRQueueLength(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取OCR队列状态失败"})
		return
	}

	status := gin.H{
		"question_queue_length":  questionQueueStatus["queue_length"],
		"image_ocr_queue_length": imageOCRQueueLength,
	}

	c.JSON(http.StatusOK, status)
}

// ProcessPendingQuestions 处理队列中的题目
func ProcessPendingQuestions(ctx context.Context) {
	fmt.Println("开始处理队列中的题目")
	for {
		fmt.Println("等待新的题目...")
		question, err := queue.BLPOPQuestion(ctx)
		if err != nil {
			if err == redis.Nil {
				// 队列为空，继续等待
				continue
			}
			fmt.Printf("从队列中获取题目失败: %v\n", err)
			// 添加一个短暂的睡眠以防止在错误情况下过度循环
			time.Sleep(time.Second)
			continue
		}

		fmt.Printf("开始处理题目: %s\n", question.ID)
		err = processQuestion(question)
		if err != nil {
			fmt.Printf("处理题目 %s 失败: %v\n", question.ID, err)
		} else {
			fmt.Printf("成功处理题目: %s\n", question.ID)
		}
	}
}

func processQuestion(question *queue.QuestionItem) error {
	// 计算内容的哈希值
	hasher := sha256.New()
	hasher.Write([]byte(question.Content))
	contentHash := hex.EncodeToString(hasher.Sum(nil))

	// 检查是否已存在相同的哈希值
	existingQuestion, err := dal.Q.Question.GetQuestionByHash(contentHash)
	if err == nil {
		// 已存在相同的题目，跳过处理
		return fmt.Errorf("题目已存在,ID: %d", existingQuestion.ID)
	}

	// 处理题目类型
	questionType, err := getOrCreateQuestionType(question.QuestionType)
	if err != nil {
		return fmt.Errorf("处理题目类型失败: %v", err)
	}

	// 处理知识点
	var knowledgePoints []model.KnowledgePoint
	for _, kpName := range question.KnowledgePoints {
		kp, err := dal.Q.KnowledgePoint.GetKnowledgePointByNameAndSubject(kpName, question.SubjectID)
		if err != nil {
			// 如果知识点不存在，创建新的知识点
			newKP := model.KnowledgePoint{
				SubjectID:   question.SubjectID,
				Name:        kpName,
				Description: nil,
				IsLeaf:      true,
			}
			err = dal.Q.KnowledgePoint.CreateKnowledgePointWithSubject(
				newKP.SubjectID,
				nil, // 假设为顶级知识点
				newKP.Name,
				newKP.Description,
				newKP.IsLeaf,
			)
			if err != nil {
				return fmt.Errorf("创建知识点失败: %v", err)
			}
			// 重新获取创建的知识点
			kp, err = dal.Q.KnowledgePoint.GetKnowledgePointByNameAndSubject(kpName, question.SubjectID)
			if err != nil {
				return fmt.Errorf("获取新创建的知识点失败: %v", err)
			}
		}
		knowledgePoints = append(knowledgePoints, *kp)
	}

	// 提取内容中的图片并进行 OCR
	contentOCRText, imagePaths, err := extractAndOCRImages(question.Content)
	if err != nil {
		return fmt.Errorf("处理内容图片失败: %v", err)
	}

	var imagePath *string
	if len(imagePaths) > 0 {
		imagePath = &imagePaths[0]
	}

	// 对答案图片进行 OCR
	answerOCRText, err := performOCR(question.ImagePath)
	if err != nil {
		return fmt.Errorf("处理答案图片失败: %v", err)
	}

	// 对解析图片进行 OCR
	explanationOCRText, err := performOCR(question.Answer)
	if err != nil {
		return fmt.Errorf("处理解析图片失败: %v", err)
	}

	// 创建新题目
	newQuestion := &model.Question{
		Content:         question.Content,
		ImagePath:       imagePath,
		OcrText:         &contentOCRText,
		Answer:          answerOCRText,
		Explanation:     &explanationOCRText,
		QuestionType:    *questionType,
		Hash:            contentHash,
		KnowledgePoints: knowledgePoints,
	}

	// 创建题目
	err = dal.Q.Question.Create(newQuestion)
	if err != nil {
		return fmt.Errorf("创建题目失败: %v", err)
	}

	return nil
}

func getOrCreateQuestionType(typeName string) (*model.QuestionType, error) {
	questionType, err := dal.Q.QuestionType.GetQuestionTypeByName(typeName)
	if err != nil {
		// 如果不存在，创建新的题目类型
		newType := &model.QuestionType{Name: typeName}
		err = dal.Q.QuestionType.CreateQuestionType(newType.Name, nil)
		if err != nil {
			return nil, fmt.Errorf("创建题目类型失败: %v", err)
		}
		// 重新获取创建的题目类型
		questionType, err = dal.Q.QuestionType.GetQuestionTypeByName(typeName)
		if err != nil {
			return nil, fmt.Errorf("获取新创建的题目类型失败: %v", err)
		}
	}
	return questionType, nil
}

func performOCR(imageURL string) (string, error) {
	ocrURL := viper.GetString("OCR_URL")

	// 创建一个multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加url字段
	_ = writer.WriteField("url", imageURL)

	// 添加其他必要的字段
	_ = writer.WriteField("multi_page", "false")
	_ = writer.WriteField("render", "false")
	_ = writer.WriteField("output_format", "plain")

	// 完成multipart写入
	err := writer.Close()
	if err != nil {
		return "", fmt.Errorf("创建multipart请求失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", ocrURL, body)
	if err != nil {
		return "", fmt.Errorf("创建OCR请求失败: %v", err)
	}

	// 设置Content-Type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送OCR请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取OCR响应失败: %v", err)
	}

	// 解析JSON响应
	var ocrResponse struct {
		Result string `json:"result"`
	}
	err = json.Unmarshal(respBody, &ocrResponse)
	if err != nil {
		return "", fmt.Errorf("解析OCR响应失败: %v", err)
	}

	return ocrResponse.Result, nil
}

func extractAndOCRImages(content string) (string, []string, error) {
	re := regexp.MustCompile(`<img[^>]+src="([^">]+)"`)
	matches := re.FindAllStringSubmatch(content, -1)

	var ocrResults []string
	var imagePaths []string
	for _, match := range matches {
		if len(match) > 1 {
			imageURL := match[1]
			imagePaths = append(imagePaths, imageURL)
			ocrText, err := performOCR(imageURL)
			if err != nil {
				return "", nil, fmt.Errorf("OCR 处理失败: %v", err)
			}
			ocrResults = append(ocrResults, ocrText)
		}
	}

	return strings.Join(ocrResults, "\n"), imagePaths, nil
}