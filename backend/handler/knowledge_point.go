package handler

import (
	"net/http"
	"strconv"

	"daofa/backend/dal"
	"daofa/backend/model"

	"github.com/gin-gonic/gin"
)

// GetKnowledgePoints 获取知识点列表（支持分页）
func GetKnowledgePoints(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 计算偏移量
	offset := (page - 1) * pageSize
	// 获取知识点列表
	knowledgePoints, err := dal.Q.KnowledgePoint.Where(dal.Q.KnowledgePoint.ID.Gt(0)).Offset(offset).Limit(pageSize).Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取知识点列表失败"})
		return
	}
	// 获取总数
	total, err := dal.Q.KnowledgePoint.CountKnowledgePoints(0, nil, "", nil, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取知识点总数失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": knowledgePoints,
		"total": total,
		"page": page,
		"pageSize": pageSize,
	})
}

// CreateKnowledgePoint 创建知识
// @router /api/v1/knowledge_point/create [POST]
func CreateKnowledgePoint(c *gin.Context) {
	var knowledgePoint model.KnowledgePoint
	if err := c.ShouldBindJSON(&knowledgePoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 1, "msg": err.Error()})
		return
	}

	err := dal.Q.KnowledgePoint.Create(&knowledgePoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "创建知识点失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": 0, "msg": "知识点创建成功", "data": knowledgePoint})
}

// UpdateKnowledgePoint 更新知识点信息
// @router /api/v1/knowledge_point/update/:id [PUT]
func UpdateKnowledgePoint(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var knowledgePoint model.KnowledgePoint
	if err := c.ShouldBindJSON(&knowledgePoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 1, "msg": err.Error()})
		return
	}

	knowledgePoint.ID = int32(id)
	_, err := dal.Q.KnowledgePoint.Where(dal.Q.KnowledgePoint.ID.Eq(knowledgePoint.ID)).Updates(knowledgePoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "更新知识点失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "知识点更新成功", "data": knowledgePoint})
}

// DeleteKnowledgePoint 删除知识点
// @router /api/v1/knowledge_point/delete/:id [DELETE]
func DeleteKnowledgePoint(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := dal.Q.KnowledgePoint.Where(dal.Q.KnowledgePoint.ID.Eq(int32(id))).Delete()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "删除知识点失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "知识点删除成功"})
}
