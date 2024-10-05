package handler

import (
	"net/http"
	"strconv"

	"daofa/backend/dal"

	"github.com/gin-gonic/gin"
)

// GetQuestionTypes 获取所有题目类型（支持分页）
func GetQuestionTypes(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// 确保 page 的值是有效的
	if page < 1 {
		page = 1
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 查询题目类型列表
	questionTypes, err := dal.Q.QuestionType.Offset(offset).Limit(pageSize).Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取习题类型列表失败"})
		return
	}

	// 获取总数
	total, err := dal.Q.QuestionType.Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取习题类型总数失败"})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"items":    questionTypes,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// CreateQuestionType 创建新的题目类型
func CreateQuestionType(c *gin.Context) {
	var input struct {
		Name        string  `json:"name" binding:"required"`
		Description *string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := dal.Q.QuestionType.CreateQuestionType(input.Name, input.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建习题类型失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "习题类型创建成功"})
}

// ListQuestionTypes 列出所有题目类型（分页）
func ListQuestionTypes(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// 计算偏移量
	offset := (page - 1) * pageSize
	// 查询题目类型列表
	questionTypes, err := dal.Q.QuestionType.Where(dal.Q.QuestionType.ID.Gt(0)).Offset(offset).Limit(pageSize).Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取习题类型列表失败"})
		return
	}

	// 获取总数
	total, err := dal.Q.QuestionType.Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取习题类型总数失败"})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"items":    questionTypes,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetQuestionType 获取单个题目类型
func GetQuestionType(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	questionType, err := dal.Q.QuestionType.GetQuestionTypeByID(int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "习题类型不存在"})
		return
	}

	c.JSON(http.StatusOK, questionType)
}

// UpdateQuestionType 更新题目类型
func UpdateQuestionType(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	var input struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := dal.Q.QuestionType.UpdateQuestionType(int32(id), input.Name, input.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新习题类型失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "习题类型更新成功"})
}

// DeleteQuestionType 删除题目类型
func DeleteQuestionType(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	_, err := dal.Q.QuestionType.Where(dal.Q.QuestionType.ID.Eq(int32(id))).Delete()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除习题类型失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "习题类型删除成功"})
}
