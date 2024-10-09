package handler

import (
	"daofa/backend/queue"
	"github.com/gin-gonic/gin"
	"net/http"
)

// EnqueueImageOCR 处理来自Chrome插件的OCR请求
func EnqueueImageOCR(c *gin.Context) {
	var request struct {
		ImageURL string `json:"imageUrl" binding:"required"`
		Cookie   string `json:"cookie"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 创建OCR任务项
	task := queue.ImageOCRTask{
		ImageURL: request.ImageURL,
		Cookie:   request.Cookie,
		Referer:  c.GetHeader("Referer"),
	}

	// 将任务加入队列
	err := queue.EnqueueImageOCRTask(c, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法将任务加入队列"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "图片已成功加入OCR处理队列"})
}