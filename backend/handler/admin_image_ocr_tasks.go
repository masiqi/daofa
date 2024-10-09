package handler

import (
	"daofa/backend/dal"
	"daofa/backend/middleware"
	"daofa/backend/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListImageOCRTasks 列出所有OCR任务
func ListImageOCRTasks(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	tasks, err := dal.Q.ImageOcrTask.ListImageOCRTasksWithPagination(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取OCR任务列表失败"})
		return
	}

	count, err := dal.Q.ImageOcrTask.CountImageOCRTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取OCR任务总数失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"total": count,
	})
}

// GetImageOCRTask 获取单个OCR任务
func GetImageOCRTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	task, err := dal.Q.ImageOcrTask.GetImageOCRTaskByID(int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "OCR任务不存在"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// CreateImageOCRTask 创建新的OCR任务
func CreateImageOCRTask(c *gin.Context) {
	var task model.ImageOcrTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	err := dal.Q.ImageOcrTask.CreateImageOCRTask(task.ImageURL, task.Cookie, task.Referer, task.LocalFilePath, task.OcrResult, *task.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建OCR任务失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "OCR任务创建成功"})
}

// UpdateImageOCRTask 更新OCR任务
func UpdateImageOCRTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var task model.ImageOcrTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	err = dal.Q.ImageOcrTask.UpdateImageOCRTask(int32(id), task.ImageURL, task.Cookie, task.Referer, task.LocalFilePath, task.OcrResult, *task.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新OCR任务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OCR任务更新成功"})
}

// DeleteImageOCRTask 删除OCR任务
func DeleteImageOCRTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	err = dal.Q.ImageOcrTask.DeleteImageOCRTask(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除OCR任务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OCR任务删除成功"})
}

// SearchImageOCRTasks 搜索OCR任务
func SearchImageOCRTasks(c *gin.Context) {
	status := c.Query("status")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	tasks, err := dal.Q.ImageOcrTask.SearchImageOCRTasks(status, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索OCR任务失败"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// RegisterAdminImageOCRTasksRoutes 注册管理员OCR任务相关路由
func RegisterAdminImageOCRTasksRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	admin.Use(middleware.JWTAuth()) // 使用JWT中间件进行身份验证

	admin.GET("/image-ocr-tasks", ListImageOCRTasks)
	admin.GET("/image-ocr-tasks/:id", GetImageOCRTask)
	admin.POST("/image-ocr-tasks", CreateImageOCRTask)
	admin.PUT("/image-ocr-tasks/:id", UpdateImageOCRTask)
	admin.DELETE("/image-ocr-tasks/:id", DeleteImageOCRTask)
	admin.GET("/image-ocr-tasks/search", SearchImageOCRTasks)
}