package handler

import (
	"context"
	"net/http"
	"strconv"

	"daofa/backend/dal"
	"daofa/backend/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetAdmins 获取所有管理员（支持分页）
// @router /api/v1/admin/list [GET]
func GetAdmins(ctx context.Context, c *gin.Context) {
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
func CreateAdmin(ctx context.Context, c *gin.Context) {
	var admin model.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 1, "msg": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "密码加密败"})
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
func UpdateAdmin(ctx context.Context, c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
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

	err := dal.Q.Admin.UpdateAdmin(uint(id), admin.Username, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "更新管理员失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "管理员更新成功"})
}

// DeleteAdmin 除管理员
// @router /api/v1/admin/delete/:id [DELETE]
func DeleteAdmin(ctx context.Context, c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := dal.Q.Admin.DeleteAdmin(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "删除管理员失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "管理员删除成功"})
}

// GetAdminByID 根据ID获取管理员
// @router /api/v1/admin/:id [GET]
func GetAdminByID(ctx context.Context, c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	admin, err := dal.Q.Admin.GetAdminByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "获取管理员信息失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "", "data": admin})
}

func AdminLogin(c *gin.Context) {
	var loginInfo model.LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 1, "msg": err.Error()})
		return
	}

	admin, err := dal.Q.Admin.Login(loginInfo.Username, loginInfo.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": 1, "msg": "用户名或密码错误"})
		return
	}

	// 生成 JWT
	token, err := generateJWT(admin.Username) // 确保这里包含用户名
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "生成 token 失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "登录成功", "data": gin.H{"token": token, "username": admin.Username}})
}
