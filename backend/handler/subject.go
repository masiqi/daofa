package handler

import (
	"net/http"
	"strconv"

	"daofa/backend/dal"
	"daofa/backend/model"

	"github.com/gin-gonic/gin"
)

// GetSubjects 获取所有科目(支持分页)
func GetSubjects(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 计算偏移量
	offset := (page - 1) * pageSize
	// 查询科目列表
	subjects, err := dal.Q.Subject.Where(dal.Q.Subject.ID.Gt(0)).Offset(offset).Limit(pageSize).Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取学科列表失败"})
		return
	}

	// 获取总数
	total, err := dal.Q.Subject.Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取学科总数失败"})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"items": subjects,
		"total": total,
		"page": page,
		"pageSize": pageSize,
	})
}

// CreateSubject 创建新科目
func CreateSubject(c *gin.Context) {
	var subject model.Subject
	// 绑定JSON数据到subject结构体
	if err := c.ShouldBindJSON(&subject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 1, "msg": err.Error()})
		return
	}

	// 创建新科目
	err := dal.Q.Subject.Create(&subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "创建科目失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": 0, "msg": "科目创建成功", "data": subject})
}

// UpdateSubject 更新科目信息
func UpdateSubject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var subject model.Subject
	// 绑定JSON数据到subject结构体
	if err := c.ShouldBindJSON(&subject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 1, "msg": err.Error()})
		return
	}

	subject.ID = int32(id)
	// 更新科目信息
	_, err := dal.Q.Subject.Where(dal.Q.Subject.ID.Eq(subject.ID)).Updates(subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "更新科目失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "科目更新成功", "data": subject})
}

// DeleteSubject 删除科目
func DeleteSubject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	// 删除指定ID的科目
	_, err := dal.Q.Subject.Where(dal.Q.Subject.ID.Eq(int32(id))).Delete()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "删除科目失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "科目删除成功"})
}
