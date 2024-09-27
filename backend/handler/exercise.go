package handler

import (
	"context"
	"net/http"
	"strconv"

	"daofa/backend/dal"
	"daofa/backend/model"

	"github.com/gin-gonic/gin"
)

// GetExercises 获取所有练习
// @router /api/v1/exercise/list [GET]
func GetExercises(ctx context.Context, c *gin.Context) {
	exercises, err := dal.Q.ExerciseMaterial.WithContext(ctx).Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "获取练习列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "",
		"data": gin.H{
			"items": exercises,
			"total": len(exercises),
		},
	})
}

// CreateExercise 创建练习
// @router /api/v1/exercise/create [POST]
func CreateExercise(ctx context.Context, c *gin.Context) {
	var exercise model.ExerciseMaterial
	if err := c.ShouldBindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 1, "msg": err.Error()})
		return
	}

	err := dal.Q.ExerciseMaterial.WithContext(ctx).Create(&exercise)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "创建练习失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": 0, "msg": "练习创建成功", "data": exercise})
}

// UpdateExercise 更新练习信息
// @router /api/v1/exercise/update/:id [PUT]
func UpdateExercise(ctx context.Context, c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var exercise model.ExerciseMaterial
	if err := c.ShouldBindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 1, "msg": err.Error()})
		return
	}

	exercise.ID = int32(id)
	_, err := dal.Q.ExerciseMaterial.WithContext(ctx).Where(dal.Q.ExerciseMaterial.ID.Eq(exercise.ID)).Updates(exercise)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "更新练习失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "练习更新成功", "data": exercise})
}

// DeleteExercise 删除练习
// @router /api/v1/exercise/delete/:id [DELETE]
func DeleteExercise(ctx context.Context, c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := dal.Q.ExerciseMaterial.WithContext(ctx).Where(dal.Q.ExerciseMaterial.ID.Eq(int32(id))).Delete()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "删除练习失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "练习删除成功"})
}
