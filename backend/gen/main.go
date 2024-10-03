package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
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
	//     {{if level != 0}}AND level=@level{{end}}
	//   {{end}}
	// ORDER BY id LIMIT @limit OFFSET @offset
	SearchKnowledgePoints(subjectID int32, parentID *int32, name string, isLeaf *bool, level int32, offset, limit int) ([]*gen.T, error)

	// SELECT COUNT(*) FROM @@table
	//   {{where}}
	//     {{if subjectID != 0}}subject_id=@subjectID{{end}}
	//     {{if parentID != nil}}AND parent_id=@parentID{{end}}
	//     {{if name != ""}}AND name LIKE CONCAT('%', @name, '%'){{end}}
	//     {{if isLeaf != nil}}AND is_leaf=@isLeaf{{end}}
	//     {{if level != 0}}AND level=@level{{end}}
	//   {{end}}
	CountKnowledgePoints(subjectID int32, parentID *int32, name string, isLeaf *bool, level int32) (int64, error)

	// INSERT INTO @@table (subject_id, parent_id, name, description, is_leaf, level) 
	// VALUES (@subjectID, @parentID, @name, @description, @isLeaf, @level)
	CreateKnowledgePoint(subjectID int32, parentID *int32, name string, description *string, isLeaf bool, level int32) error

	// UPDATE @@table SET
	//   {{if subjectID != 0}}subject_id=@subjectID,{{end}}
	//   {{if parentID != nil}}parent_id=@parentID,{{end}}
	//   {{if name != ""}}name=@name,{{end}}
	//   {{if description != nil}}description=@description,{{end}}
	//   {{if isLeaf != nil}}is_leaf=@isLeaf,{{end}}
	//   {{if level != 0}}level=@level{{end}}
	// WHERE id=@id
	UpdateKnowledgePoint(id int32, subjectID int32, parentID *int32, name string, description *string, isLeaf *bool, level int32) error

	// DELETE FROM @@table WHERE id=@id
	DeleteKnowledgePoint(id int32) error

	// SELECT * FROM @@table WHERE name = @name LIMIT 1
	GetKnowledgePointByName(name string) (*gen.T, error)
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

	// SELECT * FROM @@table WHERE question_id=@questionID
	GetKnowledgePointsByQuestionID(questionID int32) ([]*gen.T, error)

	// SELECT * FROM @@table WHERE knowledge_point_id=@knowledgePointID
	GetQuestionsByKnowledgePointID(knowledgePointID int32) ([]*gen.T, error)

	// SELECT COUNT(*) FROM @@table WHERE question_id=@questionID
	CountKnowledgePointsByQuestionID(questionID int32) (int64, error)

	// SELECT COUNT(*) FROM @@table WHERE knowledge_point_id=@knowledgePointID
	CountQuestionsByKnowledgePointID(knowledgePointID int32) (int64, error)
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

	// 为所有表生成模型和查询文件
	g.ApplyBasic(
		g.GenerateModel("admin"),
		g.GenerateModel("subject"),
		g.GenerateModel("knowledge_point"),
		g.GenerateModel("question_type"),
		g.GenerateModel("question"),
		g.GenerateModel("question_knowledge_point"),
	)

	// 应用自定义查询接口
	g.ApplyInterface(func(AdminQuerier) {}, g.GenerateModel("admin"))
	g.ApplyInterface(func(SubjectQuerier) {}, g.GenerateModel("subject"))
	g.ApplyInterface(func(KnowledgePointQuerier) {}, g.GenerateModel("knowledge_point"))
	g.ApplyInterface(func(QuestionTypeQuerier) {}, g.GenerateModel("question_type"))
	g.ApplyInterface(func(QuestionQuerier) {}, g.GenerateModel("question"))
	g.ApplyInterface(func(QuestionKnowledgePointQuerier) {}, g.GenerateModel("question_knowledge_point"))

	// 生成代码
	g.Execute()
}