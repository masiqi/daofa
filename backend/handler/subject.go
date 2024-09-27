package handler

import (
	"context"
	"net/http"
	"strconv"

	"daofa/backend/dal"
	"daofa/backend/model"

	"github.com/gin-gonic/gin"
)

// GetSubjects 获取所有科目
// @router /api/v1/subject/list [GET]
func GetSubjects(ctx context.Context, c *gin.Context) {
	subjects, err := dal.Q.Subject.WithContext(ctx).Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "获取科目列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "",
		"data": gin.H{
			"items": subjects,
			"total": len(subjects),
		},
	})
}

// CreateSubject 创建科目
// @router /api/v1/subject/create [POST]
func CreateSubject(ctx context.Context, c *gin.Context) {
	var subject model.Subject
	if err := c.ShouldBindJSON(&subject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 1, "msg": err.Error()})
		return
	}

	err := dal.Q.Subject.WithContext(ctx).Create(&subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "创建科目失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": 0, "msg": "科目创建成功", "data": subject})
}

// UpdateSubject 更新科目信息
// @router /api/v1/subject/update/:id [PUT]
func UpdateSubject(ctx context.Context, c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var subject model.Subject
	if err := c.ShouldBindJSON(&subject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 1, "msg": err.Error()})
		return
	}

	subject.ID = int32(id)
	_, err := dal.Q.Subject.WithContext(ctx).Where(dal.Q.Subject.ID.Eq(subject.ID)).Updates(subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "更新科目失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "科目更新成功", "data": subject})
}

// DeleteSubject 删除科目
// @router /api/v1/subject/delete/:id [DELETE]
func DeleteSubject(ctx context.Context, c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := dal.Q.Subject.WithContext(ctx).Where(dal.Q.Subject.ID.Eq(int32(id))).Delete()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "删除科目失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "科目删除成功"})
}
