package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// AdminQuerier defines the query interface for Admin table
type AdminQuerier interface {
	// SELECT * FROM @@table WHERE username=@username AND password=@password LIMIT 1
	Login(username, password string) (*gen.T, error)

	// INSERT INTO @@table (username, password) VALUES (@username, @password)
	CreateAdmin(username, password string) error

	// UPDATE @@table SET
	//   {{if username != ""}}username=@username,{{end}}
	//   {{if password != ""}}password=@password,{{end}}
	// WHERE id=@id
	UpdateAdmin(id uint, username, password string) error

	// DELETE FROM @@table WHERE id=@id
	DeleteAdmin(id uint) error

	// SELECT * FROM @@table LIMIT @limit OFFSET @offset
	ListAdminsWithPagination(offset, limit int) ([]*gen.T, error)

	// SELECT COUNT(*) FROM @@table
	CountAdmins() (int64, error)

	// SELECT * FROM @@table WHERE id=@id LIMIT 1
	GetAdminByID(id uint) (*gen.T, error)

	// SELECT * FROM @@table
	//   {{where}}
	//     {{if username != ""}}username LIKE CONCAT('%', @username, '%'){{end}}
	//   {{end}}
	// LIMIT @limit OFFSET @offset
	SearchAdmins(username string, offset, limit int) ([]*gen.T, error)
}

// SubjectQuerier defines the query interface for Subject table
type SubjectQuerier interface {
	// INSERT INTO @@table (name, description) VALUES (@name, @description)
	CreateSubject(name string, description *string) error

	// UPDATE @@table SET
	//   {{if name != ""}}name=@name,{{end}}
	//   {{if description != nil}}description=@description,{{end}}
	// WHERE id=@id
	UpdateSubject(id uint, name string, description *string) error

	// DELETE FROM @@table WHERE id=@id
	DeleteSubject(id uint) error

	// SELECT * FROM @@table LIMIT @limit OFFSET @offset
	ListSubjectsWithPagination(offset, limit int) ([]*gen.T, error)

	// SELECT COUNT(*) FROM @@table
	CountSubjects() (int64, error)

	// SELECT * FROM @@table WHERE id=@id LIMIT 1
	GetSubjectByID(id uint) (*gen.T, error)

	// SELECT * FROM @@table
	//   {{where}}
	//     {{if name != ""}}name LIKE CONCAT('%', @name, '%'){{end}}
	//   {{end}}
	// LIMIT @limit OFFSET @offset
	SearchSubjects(name string, offset, limit int) ([]*gen.T, error)
}

// KnowledgePointQuerier defines the query interface for KnowledgePoint table
type KnowledgePointQuerier interface {
	// INSERT INTO @@table (subject_id, parent_id, name, description, is_leaf) VALUES (@subjectID, @parentID, @name, @description, @isLeaf)
	CreateKnowledgePoint(subjectID uint, parentID *uint, name string, description *string, isLeaf bool) error

	// UPDATE @@table SET
	//   {{if subjectID != 0}}subject_id=@subjectID,{{end}}
	//   {{if parentID != nil}}parent_id=@parentID,{{end}}
	//   {{if name != ""}}name=@name,{{end}}
	//   {{if description != nil}}description=@description,{{end}}
	//   is_leaf=@isLeaf
	// WHERE id=@id
	UpdateKnowledgePoint(id, subjectID uint, parentID *uint, name string, description *string, isLeaf bool) error

	// DELETE FROM @@table WHERE id=@id
	DeleteKnowledgePoint(id uint) error

	// SELECT kp.*, s.name as subject_name, p.name as parent_name
	// FROM @@table kp
	// LEFT JOIN subject s ON kp.subject_id = s.id
	// LEFT JOIN @@table p ON kp.parent_id = p.id
	// {{where}}
	//   {{if subjectID != 0}}kp.subject_id = @subjectID{{end}}
	// {{end}}
	// LIMIT @limit OFFSET @offset
	ListKnowledgePointsWithPagination(subjectID uint, offset, limit int) ([]*gen.T, error)

	// SELECT COUNT(*) FROM @@table WHERE subject_id=@subjectID
	CountKnowledgePoints(subjectID uint) (int64, error)

	// SELECT kp.*, s.name as subject_name, p.name as parent_name
	// FROM @@table kp
	// LEFT JOIN subject s ON kp.subject_id = s.id
	// LEFT JOIN @@table p ON kp.parent_id = p.id
	// WHERE kp.id=@id LIMIT 1
	GetKnowledgePointByID(id uint) (*gen.T, error)

	// SELECT kp.*, s.name as subject_name, p.name as parent_name
	// FROM @@table kp
	// LEFT JOIN subject s ON kp.subject_id = s.id
	// LEFT JOIN @@table p ON kp.parent_id = p.id
	// {{where}}
	//   kp.subject_id=@subjectID
	//   {{if name != ""}}AND kp.name LIKE CONCAT('%', @name, '%'){{end}}
	// {{end}}
	// LIMIT @limit OFFSET @offset
	SearchKnowledgePoints(subjectID uint, name string, offset, limit int) ([]*gen.T, error)

	// SELECT kp.*, s.name as subject_name, p.name as parent_name
	// FROM @@table kp
	// LEFT JOIN subject s ON kp.subject_id = s.id
	// LEFT JOIN @@table p ON kp.parent_id = p.id
	// WHERE kp.parent_id=@parentID
	GetChildKnowledgePoints(parentID uint) ([]*gen.T, error)
}

// ExerciseMaterialQuerier defines the query interface for ExerciseMaterial table
type ExerciseMaterialQuerier interface {
	// INSERT INTO @@table (content, image_path) VALUES (@content, @imagePath)
	CreateMaterial(content, imagePath string) error

	// UPDATE @@table SET
	//   {{if content != ""}}content=@content,{{end}}
	//   {{if imagePath != ""}}image_path=@imagePath,{{end}}
	// WHERE id=@id
	UpdateMaterial(id uint, content, imagePath string) error

	// DELETE FROM @@table WHERE id=@id
	DeleteMaterial(id uint) error

	// SELECT * FROM @@table LIMIT @limit OFFSET @offset
	ListMaterialsWithPagination(offset, limit int) ([]*gen.T, error)

	// SELECT COUNT(*) FROM @@table
	CountMaterials() (int64, error)

	// SELECT * FROM @@table WHERE id=@id LIMIT 1
	GetMaterialByID(id uint) (*gen.T, error)

	// SELECT * FROM @@table
	//   {{where}}
	//     {{if content != ""}}content LIKE CONCAT('%', @content, '%'){{end}}
	//   {{end}}
	// LIMIT @limit OFFSET @offset
	SearchMaterials(content string, offset, limit int) ([]*gen.T, error)
}

// ExerciseQuestionQuerier defines the query interface for ExerciseQuestion table
type ExerciseQuestionQuerier interface {
	// INSERT INTO @@table (material_id, question, answer, explanation) VALUES (@materialID, @question, @answer, @explanation)
	CreateQuestion(materialID uint, question, answer string, explanation *string) error

	// UPDATE @@table SET
	//   {{if question != ""}}question=@question,{{end}}
	//   {{if answer != ""}}answer=@answer,{{end}}
	//   {{if explanation != nil}}explanation=@explanation,{{end}}
	// WHERE id=@id
	UpdateQuestion(id uint, question, answer string, explanation *string) error

	// DELETE FROM @@table WHERE id=@id
	DeleteQuestion(id uint) error

	// SELECT * FROM @@table WHERE material_id=@materialID LIMIT @limit OFFSET @offset
	ListQuestionsWithPagination(materialID uint, offset, limit int) ([]*gen.T, error)

	// SELECT COUNT(*) FROM @@table WHERE material_id=@materialID
	CountQuestions(materialID uint) (int64, error)

	// SELECT * FROM @@table WHERE id=@id LIMIT 1
	GetQuestionByID(id uint) (*gen.T, error)

	// SELECT * FROM @@table
	//   {{where}}
	//     material_id=@materialID
	//     {{if question != ""}}AND question LIKE CONCAT('%', @question, '%'){{end}}
	//   {{end}}
	// LIMIT @limit OFFSET @offset
	SearchQuestions(materialID uint, question string, offset, limit int) ([]*gen.T, error)
}

func main() {
	// Connect to the database
	dsn := "daofa:123456@tcp(localhost:3306)/daofa?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("cannot connect to database: %w", err))
	}

	// Initialize the generator
	g := gen.NewGenerator(gen.Config{
		OutPath:       "../dal",
		Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	// Use the database connection
	g.UseDB(db)

	// Generate models and query files for all tables
	g.ApplyBasic(
		g.GenerateModel("admin"),
		g.GenerateModel("subject"),
		g.GenerateModel("knowledge_point"),
		g.GenerateModel("exercise_material"),
		g.GenerateModel("exercise_question"),
	)

	// Apply custom query interfaces
	g.ApplyInterface(func(AdminQuerier) {}, g.GenerateModel("admin"))
	g.ApplyInterface(func(SubjectQuerier) {}, g.GenerateModel("subject"))
	g.ApplyInterface(func(KnowledgePointQuerier) {}, g.GenerateModel("knowledge_point"))
	g.ApplyInterface(func(ExerciseMaterialQuerier) {}, g.GenerateModel("exercise_material"))
	g.ApplyInterface(func(ExerciseQuestionQuerier) {}, g.GenerateModel("exercise_question"))

	// Generate code
	g.Execute()
}
