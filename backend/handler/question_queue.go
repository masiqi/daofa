package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"daofa/backend/dal"
	"daofa/backend/model"
	"daofa/backend/queue"

	"github.com/gin-gonic/gin"
)

// EnqueueQuestions 接收JSON格式的题目并将其加入队列
func EnqueueQuestions(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取请求体"})
		return
	}

	var questions []queue.QuestionItem
	err = json.Unmarshal(body, &questions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的JSON格式"})
		return
	}

	ctx := context.Background()
	for _, question := range questions {
		err := queue.EnqueueQuestion(ctx, question)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "入队失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%d 个题目已成功加入队列", len(questions))})
}

// GetQueueStatus 获取队列状态
func GetQueueStatus(c *gin.Context) {
	ctx := context.Background()
	status, err := queue.GetQueueStatus(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取队列状态失败"})
		return
	}

	c.JSON(http.StatusOK, status)
}

// ProcessPendingQuestions 处理队列中的题目
func ProcessPendingQuestions(ctx context.Context) {
	for {
		question, err := queue.DequeueQuestion(ctx)
		if err != nil {
			// 队列为空或出现错误，继续下一次循环
			continue
		}

		err = processQuestion(question)
		if err != nil {
			fmt.Printf("处理题目失败: %v\n", err)
		}
	}
}

func processQuestion(question *queue.QuestionItem) error {
	// 处理题目类型
	questionType, err := getOrCreateQuestionType(question.QuestionType)
	if err != nil {
		return fmt.Errorf("处理题目类型失败: %v", err)
	}

	// 创建新题目
	newQuestion := &model.Question{
		Content:     question.Content,
		ImagePath:   &question.ImagePath, // 修改：使用指针
		Answer:      question.Answer,
		Explanation: &question.Explanation, // 修改：使用指针
		TypeID:      questionType.ID,       // 修改：使用 int32 类型
	}

	err = dal.Q.Question.CreateQuestion(
		newQuestion.Content,
		newQuestion.ImagePath, // 修改：直接使用 ImagePath
		nil,                   // OCR text is not provided in this case
		newQuestion.Answer,
		newQuestion.Explanation, // 修改：直接使用 Explanation
		newQuestion.TypeID,
		"", // Hash will be generated in the CreateQuestion method
	)
	if err != nil {
		return fmt.Errorf("创建题目失败: %v", err)
	}

	// 处理知识点
	for _, kpName := range question.KnowledgePoints {
		knowledgePoint, err := getOrCreateKnowledgePoint(kpName)
		if err != nil {
			fmt.Printf("处理知识点失败: %v\n", err)
			continue
		}

		err = dal.Q.QuestionKnowledgePoint.CreateQuestionKnowledgePoint(newQuestion.ID, knowledgePoint.ID)
		if err != nil {
			fmt.Printf("关联题目和知识点失败: %v\n", err)
		}
	}

	fmt.Printf("成功处理题目: %s\n", newQuestion.Content)
	return nil
}

func getOrCreateQuestionType(name string) (*model.QuestionType, error) {
	questionType, err := dal.Q.QuestionType.GetQuestionTypeByName(name)
	if err == nil {
		return questionType, nil
	}

	// 如果不存在，创建新的题目类型
	err = dal.Q.QuestionType.CreateQuestionType(name, nil)
	if err != nil {
		return nil, err
	}

	return dal.Q.QuestionType.GetQuestionTypeByName(name)
}

func getOrCreateKnowledgePoint(name string) (*model.KnowledgePoint, error) {
	knowledgePoint, err := dal.Q.KnowledgePoint.GetKnowledgePointByName(name)
	if err == nil {
		return knowledgePoint, nil
	}

	// 如果不存在，创建新的知识点
	// 注意：这里我们假设所有新创建的知识点都是叶子节点，并且暂时不关联到特定学科
	err = dal.Q.KnowledgePoint.CreateKnowledgePoint(0, nil, name, nil, true, 0) // 添加了 level 参数，暂时设为 0
	if err != nil {
		return nil, err
	}

	return dal.Q.KnowledgePoint.GetKnowledgePointByName(name)
}
