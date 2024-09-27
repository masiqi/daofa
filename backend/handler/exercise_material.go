package handler

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB(database *gorm.DB) {
	db = database
}

type ExerciseMaterial struct {
	ID        uint               `json:"id" gorm:"primaryKey"`
	Content   string             `json:"content"`
	ImagePath string             `json:"image_path"`
	Questions []ExerciseQuestion `json:"questions" gorm:"foreignKey:MaterialID"`
}

// 添加这个方法来指定表名
func (ExerciseMaterial) TableName() string {
	return "exercise_material"
}

type ExerciseQuestion struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	MaterialID  uint   `json:"material_id"`
	Question    string `json:"question"`
	Answer      string `json:"answer"`
	Explanation string `json:"explanation"`
}

// 添加这个方法来指定表名
func (ExerciseQuestion) TableName() string {
	return "exercise_question"
}

func CreateExerciseMaterial(c *gin.Context) {
	var input ExerciseMaterial
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 确保 ID 为 0，让数据库自动生成
	input.ID = 0

	// 确保所有问题的 ID 和 MaterialID 都为 0
	for i := range input.Questions {
		input.Questions[i].ID = 0
		input.Questions[i].MaterialID = 0
	}

	// 开始事务
	tx := db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}

	// 创建素材和相关的题目
	if err := tx.Create(&input).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 重新加载素材，包括关联的题目
	if err := db.Preload("Questions").First(&input, input.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload material"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func ListExerciseMaterials(c *gin.Context) {
	var materials []ExerciseMaterial
	if err := db.Preload("Questions").Find(&materials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, materials)
}

func UpdateExerciseMaterial(c *gin.Context) {
	var material ExerciseMaterial
	if err := db.Where("id = ?", c.Param("id")).First(&material).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Material not found"})
		return
	}

	var input ExerciseMaterial
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 开始事务
	tx := db.Begin()

	// 检查是否上传了新图片
	if input.ImagePath != "" && input.ImagePath != material.ImagePath {
		// 删除旧图片
		if material.ImagePath != "" {
			oldImagePath := filepath.Join(".", material.ImagePath)
			if err := os.Remove(oldImagePath); err != nil {
				log.Printf("Failed to delete old image: %v", err)
				// 不要因为删除旧图片失败就中断整个更新过程
			}
		}
		// 更新图片路径
		material.ImagePath = input.ImagePath
	}

	// 更新素材内容
	material.Content = input.Content
	if err := tx.Save(&material).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 删除旧的题目
	if err := tx.Where("material_id = ?", material.ID).Delete(&ExerciseQuestion{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建新的题目
	for _, question := range input.Questions {
		question.MaterialID = material.ID
		if err := tx.Create(&question).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 重新加载素材，包括关联的题目
	if err := db.Preload("Questions").First(&material, material.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload material"})
		return
	}

	c.JSON(http.StatusOK, material)
}

func DeleteExerciseMaterial(c *gin.Context) {
	var material ExerciseMaterial
	if err := db.First(&material, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Material not found"})
		return
	}

	// 开始事务
	tx := db.Begin()

	// 删除相关的题目
	if err := tx.Where("material_id = ?", material.ID).Delete(&ExerciseQuestion{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 删除素材
	if err := tx.Delete(&material).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 删除图片文件
	if material.ImagePath != "" {
		imagePath := filepath.Join(".", material.ImagePath)
		if err := os.Remove(imagePath); err != nil {
			log.Printf("Failed to delete image file: %v", err)
			// 不要因为删除图片失败就返回错误响应
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Material and related questions deleted"})
}

// UploadImage 函数保持不变
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成文件哈希值
	hash := md5.New()
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()
	if _, err := io.Copy(hash, f); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	hashStr := fmt.Sprintf("%x", hash.Sum(nil))

	// 生成文件路径
	dir := filepath.Join("uploads", hashStr[:2], hashStr[2:4])
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 处理文件名，去掉不安全字符
	filename := filepath.Base(file.Filename)
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ReplaceAll(filename, "..", "")
	filepath := filepath.Join(dir, filename)

	// 保存文件到服务器
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回文件路径
	c.JSON(http.StatusOK, gin.H{"file_path": "/" + filepath})
}

func GetExerciseQuestions(c *gin.Context) {
	materialID := c.Query("material_id")
	if materialID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "material_id is required"})
		return
	}

	var questions []ExerciseQuestion
	if err := db.Where("material_id = ?", materialID).Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}
