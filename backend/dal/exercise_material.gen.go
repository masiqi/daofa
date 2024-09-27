// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dal

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gen/helper"

	"gorm.io/plugin/dbresolver"

	"daofa/backend/model"
)

func newExerciseMaterial(db *gorm.DB, opts ...gen.DOOption) exerciseMaterial {
	_exerciseMaterial := exerciseMaterial{}

	_exerciseMaterial.exerciseMaterialDo.UseDB(db, opts...)
	_exerciseMaterial.exerciseMaterialDo.UseModel(&model.ExerciseMaterial{})

	tableName := _exerciseMaterial.exerciseMaterialDo.TableName()
	_exerciseMaterial.ALL = field.NewAsterisk(tableName)
	_exerciseMaterial.ID = field.NewInt32(tableName, "id")
	_exerciseMaterial.Content = field.NewString(tableName, "content")
	_exerciseMaterial.ImagePath = field.NewString(tableName, "image_path")
	_exerciseMaterial.CreatedAt = field.NewTime(tableName, "created_at")
	_exerciseMaterial.UpdatedAt = field.NewTime(tableName, "updated_at")

	_exerciseMaterial.fillFieldMap()

	return _exerciseMaterial
}

type exerciseMaterial struct {
	exerciseMaterialDo

	ALL       field.Asterisk
	ID        field.Int32
	Content   field.String
	ImagePath field.String
	CreatedAt field.Time
	UpdatedAt field.Time

	fieldMap map[string]field.Expr
}

func (e exerciseMaterial) Table(newTableName string) *exerciseMaterial {
	e.exerciseMaterialDo.UseTable(newTableName)
	return e.updateTableName(newTableName)
}

func (e exerciseMaterial) As(alias string) *exerciseMaterial {
	e.exerciseMaterialDo.DO = *(e.exerciseMaterialDo.As(alias).(*gen.DO))
	return e.updateTableName(alias)
}

func (e *exerciseMaterial) updateTableName(table string) *exerciseMaterial {
	e.ALL = field.NewAsterisk(table)
	e.ID = field.NewInt32(table, "id")
	e.Content = field.NewString(table, "content")
	e.ImagePath = field.NewString(table, "image_path")
	e.CreatedAt = field.NewTime(table, "created_at")
	e.UpdatedAt = field.NewTime(table, "updated_at")

	e.fillFieldMap()

	return e
}

func (e *exerciseMaterial) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := e.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (e *exerciseMaterial) fillFieldMap() {
	e.fieldMap = make(map[string]field.Expr, 5)
	e.fieldMap["id"] = e.ID
	e.fieldMap["content"] = e.Content
	e.fieldMap["image_path"] = e.ImagePath
	e.fieldMap["created_at"] = e.CreatedAt
	e.fieldMap["updated_at"] = e.UpdatedAt
}

func (e exerciseMaterial) clone(db *gorm.DB) exerciseMaterial {
	e.exerciseMaterialDo.ReplaceConnPool(db.Statement.ConnPool)
	return e
}

func (e exerciseMaterial) replaceDB(db *gorm.DB) exerciseMaterial {
	e.exerciseMaterialDo.ReplaceDB(db)
	return e
}

type exerciseMaterialDo struct{ gen.DO }

type IExerciseMaterialDo interface {
	gen.SubQuery
	Debug() IExerciseMaterialDo
	WithContext(ctx context.Context) IExerciseMaterialDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IExerciseMaterialDo
	WriteDB() IExerciseMaterialDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IExerciseMaterialDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IExerciseMaterialDo
	Not(conds ...gen.Condition) IExerciseMaterialDo
	Or(conds ...gen.Condition) IExerciseMaterialDo
	Select(conds ...field.Expr) IExerciseMaterialDo
	Where(conds ...gen.Condition) IExerciseMaterialDo
	Order(conds ...field.Expr) IExerciseMaterialDo
	Distinct(cols ...field.Expr) IExerciseMaterialDo
	Omit(cols ...field.Expr) IExerciseMaterialDo
	Join(table schema.Tabler, on ...field.Expr) IExerciseMaterialDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IExerciseMaterialDo
	RightJoin(table schema.Tabler, on ...field.Expr) IExerciseMaterialDo
	Group(cols ...field.Expr) IExerciseMaterialDo
	Having(conds ...gen.Condition) IExerciseMaterialDo
	Limit(limit int) IExerciseMaterialDo
	Offset(offset int) IExerciseMaterialDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IExerciseMaterialDo
	Unscoped() IExerciseMaterialDo
	Create(values ...*model.ExerciseMaterial) error
	CreateInBatches(values []*model.ExerciseMaterial, batchSize int) error
	Save(values ...*model.ExerciseMaterial) error
	First() (*model.ExerciseMaterial, error)
	Take() (*model.ExerciseMaterial, error)
	Last() (*model.ExerciseMaterial, error)
	Find() ([]*model.ExerciseMaterial, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.ExerciseMaterial, err error)
	FindInBatches(result *[]*model.ExerciseMaterial, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.ExerciseMaterial) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IExerciseMaterialDo
	Assign(attrs ...field.AssignExpr) IExerciseMaterialDo
	Joins(fields ...field.RelationField) IExerciseMaterialDo
	Preload(fields ...field.RelationField) IExerciseMaterialDo
	FirstOrInit() (*model.ExerciseMaterial, error)
	FirstOrCreate() (*model.ExerciseMaterial, error)
	FindByPage(offset int, limit int) (result []*model.ExerciseMaterial, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IExerciseMaterialDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	CreateMaterial(content string, imagePath string) (err error)
	UpdateMaterial(id uint, content string, imagePath string) (err error)
	DeleteMaterial(id uint) (err error)
	ListMaterialsWithPagination(offset int, limit int) (result []*model.ExerciseMaterial, err error)
	CountMaterials() (result int64, err error)
	GetMaterialByID(id uint) (result *model.ExerciseMaterial, err error)
	SearchMaterials(content string, offset int, limit int) (result []*model.ExerciseMaterial, err error)
}

// INSERT INTO @@table (content, image_path) VALUES (@content, @imagePath)
func (e exerciseMaterialDo) CreateMaterial(content string, imagePath string) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, content)
	params = append(params, imagePath)
	generateSQL.WriteString("INSERT INTO exercise_material (content, image_path) VALUES (?, ?) ")

	var executeSQL *gorm.DB
	executeSQL = e.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

// UPDATE @@table SET
//
//	{{if content != ""}}content=@content,{{end}}
//	{{if imagePath != ""}}image_path=@imagePath,{{end}}
//
// WHERE id=@id
func (e exerciseMaterialDo) UpdateMaterial(id uint, content string, imagePath string) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("UPDATE exercise_material SET ")
	if content != "" {
		params = append(params, content)
		generateSQL.WriteString("content=?, ")
	}
	if imagePath != "" {
		params = append(params, imagePath)
		generateSQL.WriteString("image_path=?, ")
	}
	params = append(params, id)
	generateSQL.WriteString("WHERE id=? ")

	var executeSQL *gorm.DB
	executeSQL = e.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

// DELETE FROM @@table WHERE id=@id
func (e exerciseMaterialDo) DeleteMaterial(id uint) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	generateSQL.WriteString("DELETE FROM exercise_material WHERE id=? ")

	var executeSQL *gorm.DB
	executeSQL = e.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT * FROM @@table LIMIT @limit OFFSET @offset
func (e exerciseMaterialDo) ListMaterialsWithPagination(offset int, limit int) (result []*model.ExerciseMaterial, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, limit)
	params = append(params, offset)
	generateSQL.WriteString("SELECT * FROM exercise_material LIMIT ? OFFSET ? ")

	var executeSQL *gorm.DB
	executeSQL = e.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT COUNT(*) FROM @@table
func (e exerciseMaterialDo) CountMaterials() (result int64, err error) {
	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT COUNT(*) FROM exercise_material ")

	var executeSQL *gorm.DB
	executeSQL = e.UnderlyingDB().Raw(generateSQL.String()).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT * FROM @@table WHERE id=@id LIMIT 1
func (e exerciseMaterialDo) GetMaterialByID(id uint) (result *model.ExerciseMaterial, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	generateSQL.WriteString("SELECT * FROM exercise_material WHERE id=? LIMIT 1 ")

	var executeSQL *gorm.DB
	executeSQL = e.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT * FROM @@table
//
//	{{where}}
//	  {{if content != ""}}content LIKE CONCAT('%', @content, '%'){{end}}
//	{{end}}
//
// LIMIT @limit OFFSET @offset
func (e exerciseMaterialDo) SearchMaterials(content string, offset int, limit int) (result []*model.ExerciseMaterial, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM exercise_material ")
	var whereSQL0 strings.Builder
	if content != "" {
		params = append(params, content)
		whereSQL0.WriteString("content LIKE CONCAT('%', ?, '%') ")
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)
	params = append(params, limit)
	params = append(params, offset)
	generateSQL.WriteString("LIMIT ? OFFSET ? ")

	var executeSQL *gorm.DB
	executeSQL = e.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (e exerciseMaterialDo) Debug() IExerciseMaterialDo {
	return e.withDO(e.DO.Debug())
}

func (e exerciseMaterialDo) WithContext(ctx context.Context) IExerciseMaterialDo {
	return e.withDO(e.DO.WithContext(ctx))
}

func (e exerciseMaterialDo) ReadDB() IExerciseMaterialDo {
	return e.Clauses(dbresolver.Read)
}

func (e exerciseMaterialDo) WriteDB() IExerciseMaterialDo {
	return e.Clauses(dbresolver.Write)
}

func (e exerciseMaterialDo) Session(config *gorm.Session) IExerciseMaterialDo {
	return e.withDO(e.DO.Session(config))
}

func (e exerciseMaterialDo) Clauses(conds ...clause.Expression) IExerciseMaterialDo {
	return e.withDO(e.DO.Clauses(conds...))
}

func (e exerciseMaterialDo) Returning(value interface{}, columns ...string) IExerciseMaterialDo {
	return e.withDO(e.DO.Returning(value, columns...))
}

func (e exerciseMaterialDo) Not(conds ...gen.Condition) IExerciseMaterialDo {
	return e.withDO(e.DO.Not(conds...))
}

func (e exerciseMaterialDo) Or(conds ...gen.Condition) IExerciseMaterialDo {
	return e.withDO(e.DO.Or(conds...))
}

func (e exerciseMaterialDo) Select(conds ...field.Expr) IExerciseMaterialDo {
	return e.withDO(e.DO.Select(conds...))
}

func (e exerciseMaterialDo) Where(conds ...gen.Condition) IExerciseMaterialDo {
	return e.withDO(e.DO.Where(conds...))
}

func (e exerciseMaterialDo) Order(conds ...field.Expr) IExerciseMaterialDo {
	return e.withDO(e.DO.Order(conds...))
}

func (e exerciseMaterialDo) Distinct(cols ...field.Expr) IExerciseMaterialDo {
	return e.withDO(e.DO.Distinct(cols...))
}

func (e exerciseMaterialDo) Omit(cols ...field.Expr) IExerciseMaterialDo {
	return e.withDO(e.DO.Omit(cols...))
}

func (e exerciseMaterialDo) Join(table schema.Tabler, on ...field.Expr) IExerciseMaterialDo {
	return e.withDO(e.DO.Join(table, on...))
}

func (e exerciseMaterialDo) LeftJoin(table schema.Tabler, on ...field.Expr) IExerciseMaterialDo {
	return e.withDO(e.DO.LeftJoin(table, on...))
}

func (e exerciseMaterialDo) RightJoin(table schema.Tabler, on ...field.Expr) IExerciseMaterialDo {
	return e.withDO(e.DO.RightJoin(table, on...))
}

func (e exerciseMaterialDo) Group(cols ...field.Expr) IExerciseMaterialDo {
	return e.withDO(e.DO.Group(cols...))
}

func (e exerciseMaterialDo) Having(conds ...gen.Condition) IExerciseMaterialDo {
	return e.withDO(e.DO.Having(conds...))
}

func (e exerciseMaterialDo) Limit(limit int) IExerciseMaterialDo {
	return e.withDO(e.DO.Limit(limit))
}

func (e exerciseMaterialDo) Offset(offset int) IExerciseMaterialDo {
	return e.withDO(e.DO.Offset(offset))
}

func (e exerciseMaterialDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IExerciseMaterialDo {
	return e.withDO(e.DO.Scopes(funcs...))
}

func (e exerciseMaterialDo) Unscoped() IExerciseMaterialDo {
	return e.withDO(e.DO.Unscoped())
}

func (e exerciseMaterialDo) Create(values ...*model.ExerciseMaterial) error {
	if len(values) == 0 {
		return nil
	}
	return e.DO.Create(values)
}

func (e exerciseMaterialDo) CreateInBatches(values []*model.ExerciseMaterial, batchSize int) error {
	return e.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (e exerciseMaterialDo) Save(values ...*model.ExerciseMaterial) error {
	if len(values) == 0 {
		return nil
	}
	return e.DO.Save(values)
}

func (e exerciseMaterialDo) First() (*model.ExerciseMaterial, error) {
	if result, err := e.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.ExerciseMaterial), nil
	}
}

func (e exerciseMaterialDo) Take() (*model.ExerciseMaterial, error) {
	if result, err := e.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.ExerciseMaterial), nil
	}
}

func (e exerciseMaterialDo) Last() (*model.ExerciseMaterial, error) {
	if result, err := e.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.ExerciseMaterial), nil
	}
}

func (e exerciseMaterialDo) Find() ([]*model.ExerciseMaterial, error) {
	result, err := e.DO.Find()
	return result.([]*model.ExerciseMaterial), err
}

func (e exerciseMaterialDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.ExerciseMaterial, err error) {
	buf := make([]*model.ExerciseMaterial, 0, batchSize)
	err = e.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (e exerciseMaterialDo) FindInBatches(result *[]*model.ExerciseMaterial, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return e.DO.FindInBatches(result, batchSize, fc)
}

func (e exerciseMaterialDo) Attrs(attrs ...field.AssignExpr) IExerciseMaterialDo {
	return e.withDO(e.DO.Attrs(attrs...))
}

func (e exerciseMaterialDo) Assign(attrs ...field.AssignExpr) IExerciseMaterialDo {
	return e.withDO(e.DO.Assign(attrs...))
}

func (e exerciseMaterialDo) Joins(fields ...field.RelationField) IExerciseMaterialDo {
	for _, _f := range fields {
		e = *e.withDO(e.DO.Joins(_f))
	}
	return &e
}

func (e exerciseMaterialDo) Preload(fields ...field.RelationField) IExerciseMaterialDo {
	for _, _f := range fields {
		e = *e.withDO(e.DO.Preload(_f))
	}
	return &e
}

func (e exerciseMaterialDo) FirstOrInit() (*model.ExerciseMaterial, error) {
	if result, err := e.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.ExerciseMaterial), nil
	}
}

func (e exerciseMaterialDo) FirstOrCreate() (*model.ExerciseMaterial, error) {
	if result, err := e.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.ExerciseMaterial), nil
	}
}

func (e exerciseMaterialDo) FindByPage(offset int, limit int) (result []*model.ExerciseMaterial, count int64, err error) {
	result, err = e.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = e.Offset(-1).Limit(-1).Count()
	return
}

func (e exerciseMaterialDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = e.Count()
	if err != nil {
		return
	}

	err = e.Offset(offset).Limit(limit).Scan(result)
	return
}

func (e exerciseMaterialDo) Scan(result interface{}) (err error) {
	return e.DO.Scan(result)
}

func (e exerciseMaterialDo) Delete(models ...*model.ExerciseMaterial) (result gen.ResultInfo, err error) {
	return e.DO.Delete(models)
}

func (e *exerciseMaterialDo) withDO(do gen.Dao) *exerciseMaterialDo {
	e.DO = *do.(*gen.DO)
	return e
}
