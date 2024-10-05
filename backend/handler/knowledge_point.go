package handler

import (
	"net/http"
	"strconv"

	"daofa/backend/dal"

	"github.com/gin-gonic/gin"
)

// GetKnowledgePoints 获取知识点列表（支持分页）
func GetKnowledgePoints(c *gin.Context) {
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

	// 使用新的注释方法获取知识点列表
	knowledgePoints, err := dal.Q.KnowledgePoint.
		Preload(dal.Q.KnowledgePoint.Subject).
		Offset(offset).
		Limit(pageSize).
		Find()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取知识点列表失败"})
		return
	}

	// 获取总数
	total, err := dal.Q.KnowledgePoint.Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取知识点总数失败"})
		return
	}

	// 构造返回结果
	result := make([]gin.H, len(knowledgePoints))
	for i, kp := range knowledgePoints {
		result[i] = gin.H{
			"id":          kp.ID,
			"name":        kp.Name,
			"description": kp.Description,
			"is_leaf":     kp.IsLeaf,
			"subject": gin.H{
				"id":   kp.Subject.ID,
				"name": kp.Subject.Name,
			},
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"items":    result,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// CreateKnowledgePoint 创建新知识点
func CreateKnowledgePoint(c *gin.Context) {
	var input struct {
		SubjectID   int32   `json:"subject_id" binding:"required"`
		ParentID    *int32  `json:"parent_id"`
		Name        string  `json:"name" binding:"required"`
		Description *string `json:"description"`
		IsLeaf      bool    `json:"is_leaf"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := dal.Q.KnowledgePoint.CreateKnowledgePointWithSubject(
		input.SubjectID,
		input.ParentID,
		input.Name,
		input.Description,
		input.IsLeaf,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建知识点失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "知识点创建成功"})
}

// UpdateKnowledgePoint 更新知识点
func UpdateKnowledgePoint(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	var input struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
		IsLeaf      *bool   `json:"is_leaf"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Description != nil {
		updates["description"] = input.Description
	}
	if input.IsLeaf != nil {
		updates["is_leaf"] = *input.IsLeaf
	}

	_, err := dal.Q.KnowledgePoint.Where(dal.Q.KnowledgePoint.ID.Eq(int32(id))).Updates(updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新知识点失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "知识点更新成功"})
}

// DeleteKnowledgePoint 删除知识点
func DeleteKnowledgePoint(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	_, err := dal.Q.KnowledgePoint.Where(dal.Q.KnowledgePoint.ID.Eq(int32(id))).Delete()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除知识点失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "知识点删除成功"})
}

// GetKnowledgePoint 获取单个知识点
func GetKnowledgePoint(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	knowledgePoint, err := dal.Q.KnowledgePoint.
		Preload(dal.Q.KnowledgePoint.Subject).
		Where(dal.Q.KnowledgePoint.ID.Eq(int32(id))).
		First()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "知识点不存在"})
		return
	}

	result := gin.H{
		"id":          knowledgePoint.ID,
		"name":        knowledgePoint.Name,
		"description": knowledgePoint.Description,
		"is_leaf":     knowledgePoint.IsLeaf,
		"subject_id":  knowledgePoint.SubjectID,
		"parent_id":   knowledgePoint.ParentID,
		"subject": gin.H{
			"id":   knowledgePoint.Subject.ID,
			"name": knowledgePoint.Subject.Name,
		},
	}

	c.JSON(http.StatusOK, result)
}
