package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"

	"daofa/backend/dal"
	"daofa/backend/model"

	"github.com/gin-gonic/gin"
)

// CreateQuestion 创建新的题目
func CreateQuestion(c *gin.Context) {
	var input model.Question
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 计算内容的哈希值
	hasher := sha256.New()
	hasher.Write([]byte(input.Content))
	contentHash := hex.EncodeToString(hasher.Sum(nil))
	input.Hash = contentHash

	// 检查是否已存在相同的哈希值
	existingQuestion, err := dal.Q.Question.GetQuestionByHash(contentHash)
	if err == nil {
		// 已存在相同的题目
		c.JSON(http.StatusConflict, gin.H{"error": "Duplicate question", "existing_id": existingQuestion.ID})
		return
	}

	// 创建新题目
	err = dal.Q.Question.CreateQuestion(
		input.Content,
		input.ImagePath,
		input.OcrText,
		input.Answer,
		input.Explanation,
		input.TypeID,
		input.Hash,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Question created successfully"})
}

// UpdateQuestion 更新题目
func UpdateQuestion(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	var input model.Question
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新题目
	err := dal.Q.Question.UpdateQuestion(
		int32(id),
		input.Content,
		input.ImagePath,
		input.OcrText,
		input.Answer,
		input.Explanation,
		input.TypeID,
		input.Hash,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question updated successfully"})
}

// DeleteQuestion 删除题目
func DeleteQuestion(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	err := dal.Q.Question.DeleteQuestion(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}

// GetQuestion 获取单个题目
func GetQuestion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	question, err := dal.Q.Question.GetQuestionByID(int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	c.JSON(http.StatusOK, question)
}

// ListQuestions 列出所有题目（分页）
func ListQuestions(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 计算偏移量
	offset := (page - 1) * pageSize
	// 查询题目列表
	questions, err := dal.Q.Question.Where(dal.Q.Question.ID.Gt(0)).Offset(offset).Limit(pageSize).Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取习题列表失败"})
		return
	}

	// 获取总数
	total, err := dal.Q.Question.Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取习题总数失败"})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"items": questions,
		"total": total,
		"page": page,
		"pageSize": pageSize,
	})
}

// SearchQuestions 搜索题目
func SearchQuestions(c *gin.Context) {
	content := c.Query("content")
	typeID, _ := strconv.ParseInt(c.Query("type_id"), 10, 32)
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	questions, err := dal.Q.Question.SearchQuestions(content, int32(typeID), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}

// AddQuestionKnowledgePoint 为题目添加知识点
func AddQuestionKnowledgePoint(c *gin.Context) {
	questionID, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	var input struct {
		KnowledgePointID int32 `json:"knowledge_point_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := dal.Q.QuestionKnowledgePoint.CreateQuestionKnowledgePoint(int32(questionID), input.KnowledgePointID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Knowledge point added to question successfully"})
}

// RemoveQuestionKnowledgePoint 从题目中移除知识点
func RemoveQuestionKnowledgePoint(c *gin.Context) {
	questionID, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	knowledgePointID, _ := strconv.ParseInt(c.Param("knowledge_point_id"), 10, 32)

	err := dal.Q.QuestionKnowledgePoint.DeleteQuestionKnowledgePoint(int32(questionID), int32(knowledgePointID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Knowledge point removed from question successfully"})
}

// GetQuestionKnowledgePoints 获取题目的所有知识点
func GetQuestionKnowledgePoints(c *gin.Context) {
	questionID, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	knowledgePoints, err := dal.Q.QuestionKnowledgePoint.GetKnowledgePointsByQuestionID(int32(questionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, knowledgePoints)
}