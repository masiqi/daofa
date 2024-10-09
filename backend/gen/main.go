package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

// AdminQuerier 定义 Admin 表的查询接口
type AdminQuerier interface {
	// SELECT * FROM @@table WHERE username=@username AND password=@password LIMIT 1
	Login(username, password string) (*gen.T, error)

	// INSERT INTO @@table (username, password) VALUES (@username, @password)
	CreateAdmin(username, password string) error

	// UPDATE @@table SET
	//   {{if username != ""}}username=@username,{{end}}
	//   {{if password != ""}}password=@password,{{end}}
	// WHERE id=@id
	UpdateAdmin(id int32, username, password string) error

	// DELETE FROM @@table WHERE id=@id
	DeleteAdmin(id int32) error

	// SELECT * FROM @@table LIMIT @limit OFFSET @offset
	ListAdminsWithPagination(offset, limit int) ([]*gen.T, error)

	// SELECT COUNT(*) FROM @@table
	CountAdmins() (int64, error)

	// SELECT * FROM @@table WHERE id=@id LIMIT 1
	GetAdminByID(id int32) (*gen.T, error)

	// SELECT * FROM @@table
	//   {{where}}
	//     {{if username != ""}}username LIKE CONCAT('%', @username, '%'){{end}}
	//   {{end}}
	// LIMIT @limit OFFSET @offset
	SearchAdmins(username string, offset, limit int) ([]*gen.T, error)
}

// SubjectQuerier 定义 Subject 表的查询接口
type SubjectQuerier interface {
	// INSERT INTO @@table (name, description) VALUES (@name, @description)
	CreateSubject(name string, description *string) error

	// UPDATE @@table SET
	//   {{if name != ""}}name=@name,{{end}}
	//   {{if description != nil}}description=@description,{{end}}
	// WHERE id=@id
	UpdateSubject(id int32, name string, description *string) error

	// DELETE FROM @@table WHERE id=@id
	DeleteSubject(id int32) error

	// SELECT * FROM @@table LIMIT @limit OFFSET @offset
	ListSubjectsWithPagination(offset, limit int) ([]*gen.T, error)

	// SELECT COUNT(*) FROM @@table
	CountSubjects() (int64, error)

	// SELECT * FROM @@table WHERE id=@id LIMIT 1
	GetSubjectByID(id int32) (*gen.T, error)

	// SELECT * FROM @@table
	//   {{where}}
	//     {{if name != ""}}name LIKE CONCAT('%', @name, '%'){{end}}
	//   {{end}}
	// LIMIT @limit OFFSET @offset
	SearchSubjects(name string, offset, limit int) ([]*gen.T, error)
}

// KnowledgePointQuerier defines the query interface for KnowledgePoint table
type KnowledgePointQuerier interface {
	// SELECT * FROM @@table WHERE id=@id LIMIT 1
	GetKnowledgePointByID(id int32) (*gen.T, error)

	// SELECT * FROM @@table WHERE subject_id=@subjectID ORDER BY id LIMIT @limit OFFSET @offset
	ListKnowledgePointsBySubject(subjectID int32, offset, limit int) ([]*gen.T, error)

	// SELECT * FROM @@table WHERE parent_id=@parentID ORDER BY id
	GetChildKnowledgePoints(parentID int32) ([]*gen.T, error)

	// SELECT * FROM @@table WHERE subject_id=@subjectID AND parent_id IS NULL ORDER BY id
	GetRootKnowledgePoints(subjectID int32) ([]*gen.T, error)

	// SELECT * FROM @@table
	//   {{where}}
	//     {{if subjectID != 0}}subject_id=@subjectID{{end}}
	//     {{if parentID != nil}}AND parent_id=@parentID{{end}}
	//     {{if name != ""}}AND name LIKE CONCAT('%', @name, '%'){{end}}
	//     {{if isLeaf != nil}}AND is_leaf=@isLeaf{{end}}
	//   {{end}}
	// ORDER BY id LIMIT @limit OFFSET @offset
	SearchKnowledgePoints(subjectID int32, parentID *int32, name string, isLeaf *bool, offset, limit int) ([]*gen.T, error)

	// SELECT COUNT(*) FROM @@table
	//   {{where}}
	//     {{if subjectID != 0}}subject_id=@subjectID{{end}}
	//     {{if parentID != nil}}AND parent_id=@parentID{{end}}
	//     {{if name != ""}}AND name LIKE CONCAT('%', @name, '%'){{end}}
	//     {{if isLeaf != nil}}AND is_leaf=@isLeaf{{end}}
	//   {{end}}
	CountKnowledgePoints(subjectID int32, parentID *int32, name string, isLeaf *bool) (int64, error)

	// SELECT * FROM @@table ORDER BY id LIMIT @limit OFFSET @offset
	ListKnowledgePointsWithPagination(offset, limit int) ([]*gen.T, error)

	// INSERT INTO @@table (subject_id, parent_id, name, description, is_leaf )
	// VALUES (@subjectID, @parentID, @name, @description, @isLeaf )
	CreateKnowledgePoint(subjectID int32, parentID *int32, name string, description *string, isLeaf bool) error

	// UPDATE @@table SET
	//   {{if subjectID != 0}}subject_id=@subjectID,{{end}}
	//   {{if parentID != nil}}parent_id=@parentID,{{end}}
	//   {{if name != ""}}name=@name,{{end}}
	//   {{if description != nil}}description=@description,{{end}}
	//   {{if isLeaf != nil}}is_leaf=@isLeaf,{{end}}
	// WHERE id=@id
	UpdateKnowledgePoint(id int32, subjectID int32, parentID *int32, name string, description *string, isLeaf *bool) error

	// DELETE FROM @@table WHERE id=@id
	DeleteKnowledgePoint(id int32) error

	// SELECT * FROM @@table WHERE name = @name LIMIT 1
	GetKnowledgePointByName(name string) (*gen.T, error)

	// SELECT * FROM @@table WHERE subject_id = @subjectID AND name = @name LIMIT 1
	GetKnowledgePointByNameAndSubject(name string, subjectID int32) (*gen.T, error)

	// INSERT INTO @@table (subject_id, parent_id, name, description, is_leaf)
	// VALUES (@subjectID, @parentID, @name, @description, @isLeaf)
	CreateKnowledgePointWithSubject(subjectID int32, parentID *int32, name string, description *string, isLeaf bool) error
}

// QuestionTypeQuerier 定义 QuestionType 表的查询接口
type QuestionTypeQuerier interface {
	// INSERT INTO @@table (name, description) VALUES (@name, @description)
	CreateQuestionType(name string, description *string) error

	// UPDATE @@table SET
	//   {{if name != ""}}name=@name,{{end}}
	//   {{if description != nil}}description=@description,{{end}}
	// WHERE id=@id
	UpdateQuestionType(id int32, name string, description *string) error

	// DELETE FROM @@table WHERE id=@id
	DeleteQuestionType(id int32) error

	// SELECT * FROM @@table LIMIT @limit OFFSET @offset
	ListQuestionTypesWithPagination(offset, limit int) ([]*gen.T, error)

	// SELECT COUNT(*) FROM @@table
	CountQuestionTypes() (int64, error)

	// SELECT * FROM @@table WHERE id=@id LIMIT 1
	GetQuestionTypeByID(id int32) (*gen.T, error)

	// SELECT * FROM @@table
	//   {{where}}
	//     {{if name != ""}}name LIKE CONCAT('%', @name, '%'){{end}}
	//   {{end}}
	// LIMIT @limit OFFSET @offset
	SearchQuestionTypes(name string, offset, limit int) ([]*gen.T, error)

	// SELECT * FROM @@table WHERE name = @name LIMIT 1
	GetQuestionTypeByName(name string) (*gen.T, error)
}

// QuestionQuerier 定义 Question 表的查询接口
type QuestionQuerier interface {
	// INSERT INTO @@table (content, image_path, ocr_text, answer, explanation, type_id, hash)
	// VALUES (@content, @imagePath, @ocrText, @answer, @explanation, @typeID, @hash)
	CreateQuestion(content string, imagePath *string, ocrText *string, answer string, explanation *string, typeID int32, hash string) error

	// UPDATE @@table SET
	//   {{if content != ""}}content=@content,{{end}}
	//   {{if imagePath != nil}}image_path=@imagePath,{{end}}
	//   {{if ocrText != nil}}ocr_text=@ocrText,{{end}}
	//   {{if answer != ""}}answer=@answer,{{end}}
	//   {{if explanation != nil}}explanation=@explanation,{{end}}
	//   {{if typeID != 0}}type_id=@typeID,{{end}}
	//   {{if hash != ""}}hash=@hash{{end}}
	// WHERE id=@id
	UpdateQuestion(id int32, content string, imagePath *string, ocrText *string, answer string, explanation *string, typeID int32, hash string) error

	// DELETE FROM @@table WHERE id=@id
	DeleteQuestion(id int32) error

	// SELECT * FROM @@table LIMIT @limit OFFSET @offset
	ListQuestionsWithPagination(offset, limit int) ([]*gen.T, error)

	// SELECT COUNT(*) FROM @@table
	CountQuestions() (int64, error)

	// SELECT * FROM @@table WHERE id=@id LIMIT 1
	GetQuestionByID(id int32) (*gen.T, error)

	// SELECT * FROM @@table WHERE hash=@hash LIMIT 1
	GetQuestionByHash(hash string) (*gen.T, error)

	// SELECT * FROM @@table
	//   {{where}}
	//     {{if content != ""}}content LIKE CONCAT('%', @content, '%'){{end}}
	//     {{if typeID != 0}}AND type_id = @typeID{{end}}
	//   {{end}}
	// ORDER BY id DESC
	// LIMIT @limit OFFSET @offset
	SearchQuestions(content string, typeID int32, offset, limit int) ([]*gen.T, error)
}

// QuestionKnowledgePointQuerier 定义 QuestionKnowledgePoint 表的查询接口
type QuestionKnowledgePointQuerier interface {
	// INSERT INTO @@table (question_id, knowledge_point_id) VALUES (@questionID, @knowledgePointID)
	CreateQuestionKnowledgePoint(questionID, knowledgePointID int32) error

	// DELETE FROM @@table WHERE question_id=@questionID AND knowledge_point_id=@knowledgePointID
	DeleteQuestionKnowledgePoint(questionID, knowledgePointID int32) error

	// SELECT * FROM @@table WHERE question_id=@questionID LIMIT @limit OFFSET @offset
	GetKnowledgePointsByQuestionIDWithPagination(questionID int32, offset, limit int) ([]*gen.T, error)

	// SELECT COUNT(*) FROM @@table WHERE question_id=@questionID
	CountKnowledgePointsByQuestionID(questionID int32) (int64, error)

	// SELECT * FROM @@table WHERE knowledge_point_id=@knowledgePointID LIMIT @limit OFFSET @offset
	GetQuestionsByKnowledgePointIDWithPagination(knowledgePointID int32, offset, limit int) ([]*gen.T, error)

	// SELECT COUNT(*) FROM @@table WHERE knowledge_point_id=@knowledgePointID
	CountQuestionsByKnowledgePointID(knowledgePointID int32) (int64, error)

	// SELECT k.* FROM @@table qkp
	// JOIN knowledge_point k ON qkp.knowledge_point_id = k.id
	// WHERE qkp.question_id = @questionID
	GetKnowledgePointsByQuestionID(questionID int32) ([]*gen.T, error)
}

// ImageOCRTaskQuerier 定义 ImageOCRTask 表的查询接口
type ImageOCRTaskQuerier interface {
	// INSERT INTO @@table (image_url, cookie, referer, local_file_path, ocr_result, status)
	// VALUES (@imageURL, @cookie, @referer, @localFilePath, @ocrResult, @status)
	CreateImageOCRTask(imageURL string, cookie *string, referer *string, localFilePath *string, ocrResult *string, status string) error

	// UPDATE @@table SET
	//   {{if imageURL != ""}}image_url=@imageURL,{{end}}
	//   {{if cookie != nil}}cookie=@cookie,{{end}}
	//   {{if referer != nil}}referer=@referer,{{end}}
	//   {{if localFilePath != nil}}local_file_path=@localFilePath,{{end}}
	//   {{if ocrResult != nil}}ocr_result=@ocrResult,{{end}}
	//   {{if status != ""}}status=@status,{{end}}
	// WHERE id=@id
	UpdateImageOCRTask(id int32, imageURL string, cookie *string, referer *string, localFilePath *string, ocrResult *string, status string) error

	// SELECT * FROM @@table WHERE id=@id LIMIT 1
	GetImageOCRTaskByID(id int32) (*gen.T, error)

	// SELECT * FROM @@table ORDER BY id DESC LIMIT @limit OFFSET @offset
	ListImageOCRTasksWithPagination(offset, limit int) ([]*gen.T, error)

	// SELECT COUNT(*) FROM @@table
	CountImageOCRTasks() (int64, error)

	// SELECT * FROM @@table
	//   {{where}}
	//     {{if status != ""}}status = @status{{end}}
	//   {{end}}
	// ORDER BY id DESC
	// LIMIT @limit OFFSET @offset
	SearchImageOCRTasks(status string, offset, limit int) ([]*gen.T, error)

	// DELETE FROM @@table WHERE id=@id
	DeleteImageOCRTask(id int32) error
}

func main() {
	// 连接到数据库
	dsn := "daofa:123456@tcp(localhost:3306)/daofa?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("无法连接到数据库: %w", err))
	}

	// 初始化生成器
	g := gen.NewGenerator(gen.Config{
		OutPath:       "../dal",
		Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
		FieldSignable: true,
	})

	// 使用数据库连接
	g.UseDB(db)

	// 定义模型之间的关系
	subject := g.GenerateModel("subject")

	knowledgePoint := g.GenerateModel("knowledge_point",
		gen.FieldRelate(field.BelongsTo, "Subject", subject, &field.RelateConfig{
			GORMTag: field.GormTag{"foreignKey": []string{"subject_id"}},
		}),
		gen.FieldRelate(field.HasMany, "Children", g.GenerateModel("knowledge_point"), &field.RelateConfig{
			RelateSlice: true,
			GORMTag:     field.GormTag{"foreignKey": []string{"parent_id"}},
		}),
		gen.FieldRelate(field.BelongsTo, "Parent", g.GenerateModel("knowledge_point"), &field.RelateConfig{
			RelatePointer: true,
			GORMTag:       field.GormTag{"foreignKey": []string{"parent_id"}},
		}),
	)

	questionType := g.GenerateModel("question_type")

	question := g.GenerateModel("question",
		gen.FieldRelate(field.BelongsTo, "QuestionType", questionType, &field.RelateConfig{
			GORMTag: field.GormTag{"foreignKey": []string{"type_id"}},
		}),
		gen.FieldRelate(field.Many2Many, "KnowledgePoints", knowledgePoint, &field.RelateConfig{
			RelateSlice: true,
			GORMTag:     field.GormTag{"many2many": []string{"question_knowledge_point"}},
		}),
	)

	// 为所有表生成模型和查询文件
	g.ApplyBasic(
		g.GenerateModel("admin"),
		subject,
		knowledgePoint,
		questionType,
		question,
		g.GenerateModel("question_knowledge_point"),
		g.GenerateModel("image_ocr_tasks"),
	)

	// 应用自定义查询接口
	g.ApplyInterface(func(AdminQuerier) {}, g.GenerateModel("admin"))
	g.ApplyInterface(func(SubjectQuerier) {}, subject)
	g.ApplyInterface(func(KnowledgePointQuerier) {}, knowledgePoint)
	g.ApplyInterface(func(QuestionTypeQuerier) {}, questionType)
	g.ApplyInterface(func(QuestionQuerier) {}, question)
	g.ApplyInterface(func(QuestionKnowledgePointQuerier) {}, g.GenerateModel("question_knowledge_point"))
	// 添加新的接口
	g.ApplyInterface(func(ImageOCRTaskQuerier) {}, g.GenerateModel("image_ocr_tasks"))

	// 生成代码
	g.Execute()
}
