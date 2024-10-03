package handler

import (
	"net/http"
	"strconv"

	"daofa/backend/dal"
	"daofa/backend/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetAdmins 获取所有管理员（支持分页）
// @router /api/v1/admin/list [GET]
func GetAdmins(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "10"))

	// 计算偏移量
	offset := (page - 1) * perPage

	// 获取总数
	total, err := dal.Q.Admin.CountAdmins()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "获取管理员总数失败"})
		return
	}

	// 获取分页后的管理员列表
	admins, err := dal.Q.Admin.ListAdminsWithPagination(offset, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "获取管理员列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "",
		"data": gin.H{
			"items":   admins,
			"total":   total,
			"page":    page,
			"perPage": perPage,
		},
	})
}

// CreateAdmin 创建管理员
// @router /api/v1/admin/create [POST]
func CreateAdmin(c *gin.Context) {
	var admin model.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 1, "msg": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "密码加密失败"})
		return
	}

	err = dal.Q.Admin.CreateAdmin(admin.Username, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "创建管理员失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": 0, "msg": "管理员创建成功"})
}

// UpdateAdmin 更新管理员信息
// @router /api/v1/admin/update/:id [PUT]
func UpdateAdmin(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	var admin model.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 1, "msg": err.Error()})
		return
	}

	var hashedPassword string
	if admin.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "密码加密失败"})
			return
		}
		hashedPassword = string(hashed)
	}

	err := dal.Q.Admin.UpdateAdmin(int32(id), admin.Username, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "更新管理员失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "管理员更新成功"})
}

// DeleteAdmin 删除管理员
// @router /api/v1/admin/delete/:id [DELETE]
func DeleteAdmin(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	err := dal.Q.Admin.DeleteAdmin(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "删除管理员失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "管理员删除成功"})
}

// GetAdminByID 根据ID获取管理员
// @router /api/v1/admin/:id [GET]
func GetAdminByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	admin, err := dal.Q.Admin.GetAdminByID(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "获取管理员信息失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "", "data": admin})
}
