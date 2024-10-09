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

func newImageOcrTask(db *gorm.DB, opts ...gen.DOOption) imageOcrTask {
	_imageOcrTask := imageOcrTask{}

	_imageOcrTask.imageOcrTaskDo.UseDB(db, opts...)
	_imageOcrTask.imageOcrTaskDo.UseModel(&model.ImageOcrTask{})

	tableName := _imageOcrTask.imageOcrTaskDo.TableName()
	_imageOcrTask.ALL = field.NewAsterisk(tableName)
	_imageOcrTask.ID = field.NewInt32(tableName, "id")
	_imageOcrTask.ImageURL = field.NewString(tableName, "image_url")
	_imageOcrTask.Cookie = field.NewString(tableName, "cookie")
	_imageOcrTask.Referer = field.NewString(tableName, "referer")
	_imageOcrTask.LocalFilePath = field.NewString(tableName, "local_file_path")
	_imageOcrTask.OcrResult = field.NewString(tableName, "ocr_result")
	_imageOcrTask.Status = field.NewString(tableName, "status")
	_imageOcrTask.CreatedAt = field.NewTime(tableName, "created_at")
	_imageOcrTask.UpdatedAt = field.NewTime(tableName, "updated_at")

	_imageOcrTask.fillFieldMap()

	return _imageOcrTask
}

type imageOcrTask struct {
	imageOcrTaskDo

	ALL           field.Asterisk
	ID            field.Int32
	ImageURL      field.String
	Cookie        field.String
	Referer       field.String
	LocalFilePath field.String
	OcrResult     field.String
	Status        field.String
	CreatedAt     field.Time
	UpdatedAt     field.Time

	fieldMap map[string]field.Expr
}

func (i imageOcrTask) Table(newTableName string) *imageOcrTask {
	i.imageOcrTaskDo.UseTable(newTableName)
	return i.updateTableName(newTableName)
}

func (i imageOcrTask) As(alias string) *imageOcrTask {
	i.imageOcrTaskDo.DO = *(i.imageOcrTaskDo.As(alias).(*gen.DO))
	return i.updateTableName(alias)
}

func (i *imageOcrTask) updateTableName(table string) *imageOcrTask {
	i.ALL = field.NewAsterisk(table)
	i.ID = field.NewInt32(table, "id")
	i.ImageURL = field.NewString(table, "image_url")
	i.Cookie = field.NewString(table, "cookie")
	i.Referer = field.NewString(table, "referer")
	i.LocalFilePath = field.NewString(table, "local_file_path")
	i.OcrResult = field.NewString(table, "ocr_result")
	i.Status = field.NewString(table, "status")
	i.CreatedAt = field.NewTime(table, "created_at")
	i.UpdatedAt = field.NewTime(table, "updated_at")

	i.fillFieldMap()

	return i
}

func (i *imageOcrTask) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := i.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (i *imageOcrTask) fillFieldMap() {
	i.fieldMap = make(map[string]field.Expr, 9)
	i.fieldMap["id"] = i.ID
	i.fieldMap["image_url"] = i.ImageURL
	i.fieldMap["cookie"] = i.Cookie
	i.fieldMap["referer"] = i.Referer
	i.fieldMap["local_file_path"] = i.LocalFilePath
	i.fieldMap["ocr_result"] = i.OcrResult
	i.fieldMap["status"] = i.Status
	i.fieldMap["created_at"] = i.CreatedAt
	i.fieldMap["updated_at"] = i.UpdatedAt
}

func (i imageOcrTask) clone(db *gorm.DB) imageOcrTask {
	i.imageOcrTaskDo.ReplaceConnPool(db.Statement.ConnPool)
	return i
}

func (i imageOcrTask) replaceDB(db *gorm.DB) imageOcrTask {
	i.imageOcrTaskDo.ReplaceDB(db)
	return i
}

type imageOcrTaskDo struct{ gen.DO }

type IImageOcrTaskDo interface {
	gen.SubQuery
	Debug() IImageOcrTaskDo
	WithContext(ctx context.Context) IImageOcrTaskDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IImageOcrTaskDo
	WriteDB() IImageOcrTaskDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IImageOcrTaskDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IImageOcrTaskDo
	Not(conds ...gen.Condition) IImageOcrTaskDo
	Or(conds ...gen.Condition) IImageOcrTaskDo
	Select(conds ...field.Expr) IImageOcrTaskDo
	Where(conds ...gen.Condition) IImageOcrTaskDo
	Order(conds ...field.Expr) IImageOcrTaskDo
	Distinct(cols ...field.Expr) IImageOcrTaskDo
	Omit(cols ...field.Expr) IImageOcrTaskDo
	Join(table schema.Tabler, on ...field.Expr) IImageOcrTaskDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IImageOcrTaskDo
	RightJoin(table schema.Tabler, on ...field.Expr) IImageOcrTaskDo
	Group(cols ...field.Expr) IImageOcrTaskDo
	Having(conds ...gen.Condition) IImageOcrTaskDo
	Limit(limit int) IImageOcrTaskDo
	Offset(offset int) IImageOcrTaskDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IImageOcrTaskDo
	Unscoped() IImageOcrTaskDo
	Create(values ...*model.ImageOcrTask) error
	CreateInBatches(values []*model.ImageOcrTask, batchSize int) error
	Save(values ...*model.ImageOcrTask) error
	First() (*model.ImageOcrTask, error)
	Take() (*model.ImageOcrTask, error)
	Last() (*model.ImageOcrTask, error)
	Find() ([]*model.ImageOcrTask, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.ImageOcrTask, err error)
	FindInBatches(result *[]*model.ImageOcrTask, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.ImageOcrTask) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IImageOcrTaskDo
	Assign(attrs ...field.AssignExpr) IImageOcrTaskDo
	Joins(fields ...field.RelationField) IImageOcrTaskDo
	Preload(fields ...field.RelationField) IImageOcrTaskDo
	FirstOrInit() (*model.ImageOcrTask, error)
	FirstOrCreate() (*model.ImageOcrTask, error)
	FindByPage(offset int, limit int) (result []*model.ImageOcrTask, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IImageOcrTaskDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	CreateImageOCRTask(imageURL string, cookie *string, referer *string, localFilePath *string, ocrResult *string, status string) (err error)
	UpdateImageOCRTask(id int32, imageURL string, cookie *string, referer *string, localFilePath *string, ocrResult *string, status string) (err error)
	GetImageOCRTaskByID(id int32) (result *model.ImageOcrTask, err error)
	ListImageOCRTasksWithPagination(offset int, limit int) (result []*model.ImageOcrTask, err error)
	CountImageOCRTasks() (result int64, err error)
	SearchImageOCRTasks(status string, offset int, limit int) (result []*model.ImageOcrTask, err error)
	DeleteImageOCRTask(id int32) (err error)
}

// INSERT INTO @@table (image_url, cookie, referer, local_file_path, ocr_result, status)
// VALUES (@imageURL, @cookie, @referer, @localFilePath, @ocrResult, @status)
func (i imageOcrTaskDo) CreateImageOCRTask(imageURL string, cookie *string, referer *string, localFilePath *string, ocrResult *string, status string) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, imageURL)
	params = append(params, cookie)
	params = append(params, referer)
	params = append(params, localFilePath)
	params = append(params, ocrResult)
	params = append(params, status)
	generateSQL.WriteString("INSERT INTO image_ocr_tasks (image_url, cookie, referer, local_file_path, ocr_result, status) VALUES (?, ?, ?, ?, ?, ?) ")

	var executeSQL *gorm.DB
	executeSQL = i.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

// UPDATE @@table SET
//
//	{{if imageURL != ""}}image_url=@imageURL,{{end}}
//	{{if cookie != nil}}cookie=@cookie,{{end}}
//	{{if referer != nil}}referer=@referer,{{end}}
//	{{if localFilePath != nil}}local_file_path=@localFilePath,{{end}}
//	{{if ocrResult != nil}}ocr_result=@ocrResult,{{end}}
//	{{if status != ""}}status=@status,{{end}}
//
// WHERE id=@id
func (i imageOcrTaskDo) UpdateImageOCRTask(id int32, imageURL string, cookie *string, referer *string, localFilePath *string, ocrResult *string, status string) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("UPDATE image_ocr_tasks SET ")
	if imageURL != "" {
		params = append(params, imageURL)
		generateSQL.WriteString("image_url=?, ")
	}
	if cookie != nil {
		params = append(params, cookie)
		generateSQL.WriteString("cookie=?, ")
	}
	if referer != nil {
		params = append(params, referer)
		generateSQL.WriteString("referer=?, ")
	}
	if localFilePath != nil {
		params = append(params, localFilePath)
		generateSQL.WriteString("local_file_path=?, ")
	}
	if ocrResult != nil {
		params = append(params, ocrResult)
		generateSQL.WriteString("ocr_result=?, ")
	}
	if status != "" {
		params = append(params, status)
		generateSQL.WriteString("status=?, ")
	}
	params = append(params, id)
	generateSQL.WriteString("WHERE id=? ")

	var executeSQL *gorm.DB
	executeSQL = i.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT * FROM @@table WHERE id=@id LIMIT 1
func (i imageOcrTaskDo) GetImageOCRTaskByID(id int32) (result *model.ImageOcrTask, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	generateSQL.WriteString("SELECT * FROM image_ocr_tasks WHERE id=? LIMIT 1 ")

	var executeSQL *gorm.DB
	executeSQL = i.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT * FROM @@table ORDER BY id DESC LIMIT @limit OFFSET @offset
func (i imageOcrTaskDo) ListImageOCRTasksWithPagination(offset int, limit int) (result []*model.ImageOcrTask, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, limit)
	params = append(params, offset)
	generateSQL.WriteString("SELECT * FROM image_ocr_tasks ORDER BY id DESC LIMIT ? OFFSET ? ")

	var executeSQL *gorm.DB
	executeSQL = i.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT COUNT(*) FROM @@table
func (i imageOcrTaskDo) CountImageOCRTasks() (result int64, err error) {
	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT COUNT(*) FROM image_ocr_tasks ")

	var executeSQL *gorm.DB
	executeSQL = i.UnderlyingDB().Raw(generateSQL.String()).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT * FROM @@table
//
//	{{where}}
//	  {{if status != ""}}status = @status{{end}}
//	{{end}}
//
// ORDER BY id DESC
// LIMIT @limit OFFSET @offset
func (i imageOcrTaskDo) SearchImageOCRTasks(status string, offset int, limit int) (result []*model.ImageOcrTask, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM image_ocr_tasks ")
	var whereSQL0 strings.Builder
	if status != "" {
		params = append(params, status)
		whereSQL0.WriteString("status = ? ")
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)
	params = append(params, limit)
	params = append(params, offset)
	generateSQL.WriteString("ORDER BY id DESC LIMIT ? OFFSET ? ")

	var executeSQL *gorm.DB
	executeSQL = i.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// DELETE FROM @@table WHERE id=@id
func (i imageOcrTaskDo) DeleteImageOCRTask(id int32) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	generateSQL.WriteString("DELETE FROM image_ocr_tasks WHERE id=? ")

	var executeSQL *gorm.DB
	executeSQL = i.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (i imageOcrTaskDo) Debug() IImageOcrTaskDo {
	return i.withDO(i.DO.Debug())
}

func (i imageOcrTaskDo) WithContext(ctx context.Context) IImageOcrTaskDo {
	return i.withDO(i.DO.WithContext(ctx))
}

func (i imageOcrTaskDo) ReadDB() IImageOcrTaskDo {
	return i.Clauses(dbresolver.Read)
}

func (i imageOcrTaskDo) WriteDB() IImageOcrTaskDo {
	return i.Clauses(dbresolver.Write)
}

func (i imageOcrTaskDo) Session(config *gorm.Session) IImageOcrTaskDo {
	return i.withDO(i.DO.Session(config))
}

func (i imageOcrTaskDo) Clauses(conds ...clause.Expression) IImageOcrTaskDo {
	return i.withDO(i.DO.Clauses(conds...))
}

func (i imageOcrTaskDo) Returning(value interface{}, columns ...string) IImageOcrTaskDo {
	return i.withDO(i.DO.Returning(value, columns...))
}

func (i imageOcrTaskDo) Not(conds ...gen.Condition) IImageOcrTaskDo {
	return i.withDO(i.DO.Not(conds...))
}

func (i imageOcrTaskDo) Or(conds ...gen.Condition) IImageOcrTaskDo {
	return i.withDO(i.DO.Or(conds...))
}

func (i imageOcrTaskDo) Select(conds ...field.Expr) IImageOcrTaskDo {
	return i.withDO(i.DO.Select(conds...))
}

func (i imageOcrTaskDo) Where(conds ...gen.Condition) IImageOcrTaskDo {
	return i.withDO(i.DO.Where(conds...))
}

func (i imageOcrTaskDo) Order(conds ...field.Expr) IImageOcrTaskDo {
	return i.withDO(i.DO.Order(conds...))
}

func (i imageOcrTaskDo) Distinct(cols ...field.Expr) IImageOcrTaskDo {
	return i.withDO(i.DO.Distinct(cols...))
}

func (i imageOcrTaskDo) Omit(cols ...field.Expr) IImageOcrTaskDo {
	return i.withDO(i.DO.Omit(cols...))
}

func (i imageOcrTaskDo) Join(table schema.Tabler, on ...field.Expr) IImageOcrTaskDo {
	return i.withDO(i.DO.Join(table, on...))
}

func (i imageOcrTaskDo) LeftJoin(table schema.Tabler, on ...field.Expr) IImageOcrTaskDo {
	return i.withDO(i.DO.LeftJoin(table, on...))
}

func (i imageOcrTaskDo) RightJoin(table schema.Tabler, on ...field.Expr) IImageOcrTaskDo {
	return i.withDO(i.DO.RightJoin(table, on...))
}

func (i imageOcrTaskDo) Group(cols ...field.Expr) IImageOcrTaskDo {
	return i.withDO(i.DO.Group(cols...))
}

func (i imageOcrTaskDo) Having(conds ...gen.Condition) IImageOcrTaskDo {
	return i.withDO(i.DO.Having(conds...))
}

func (i imageOcrTaskDo) Limit(limit int) IImageOcrTaskDo {
	return i.withDO(i.DO.Limit(limit))
}

func (i imageOcrTaskDo) Offset(offset int) IImageOcrTaskDo {
	return i.withDO(i.DO.Offset(offset))
}

func (i imageOcrTaskDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IImageOcrTaskDo {
	return i.withDO(i.DO.Scopes(funcs...))
}

func (i imageOcrTaskDo) Unscoped() IImageOcrTaskDo {
	return i.withDO(i.DO.Unscoped())
}

func (i imageOcrTaskDo) Create(values ...*model.ImageOcrTask) error {
	if len(values) == 0 {
		return nil
	}
	return i.DO.Create(values)
}

func (i imageOcrTaskDo) CreateInBatches(values []*model.ImageOcrTask, batchSize int) error {
	return i.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (i imageOcrTaskDo) Save(values ...*model.ImageOcrTask) error {
	if len(values) == 0 {
		return nil
	}
	return i.DO.Save(values)
}

func (i imageOcrTaskDo) First() (*model.ImageOcrTask, error) {
	if result, err := i.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.ImageOcrTask), nil
	}
}

func (i imageOcrTaskDo) Take() (*model.ImageOcrTask, error) {
	if result, err := i.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.ImageOcrTask), nil
	}
}

func (i imageOcrTaskDo) Last() (*model.ImageOcrTask, error) {
	if result, err := i.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.ImageOcrTask), nil
	}
}

func (i imageOcrTaskDo) Find() ([]*model.ImageOcrTask, error) {
	result, err := i.DO.Find()
	return result.([]*model.ImageOcrTask), err
}

func (i imageOcrTaskDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.ImageOcrTask, err error) {
	buf := make([]*model.ImageOcrTask, 0, batchSize)
	err = i.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (i imageOcrTaskDo) FindInBatches(result *[]*model.ImageOcrTask, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return i.DO.FindInBatches(result, batchSize, fc)
}

func (i imageOcrTaskDo) Attrs(attrs ...field.AssignExpr) IImageOcrTaskDo {
	return i.withDO(i.DO.Attrs(attrs...))
}

func (i imageOcrTaskDo) Assign(attrs ...field.AssignExpr) IImageOcrTaskDo {
	return i.withDO(i.DO.Assign(attrs...))
}

func (i imageOcrTaskDo) Joins(fields ...field.RelationField) IImageOcrTaskDo {
	for _, _f := range fields {
		i = *i.withDO(i.DO.Joins(_f))
	}
	return &i
}

func (i imageOcrTaskDo) Preload(fields ...field.RelationField) IImageOcrTaskDo {
	for _, _f := range fields {
		i = *i.withDO(i.DO.Preload(_f))
	}
	return &i
}

func (i imageOcrTaskDo) FirstOrInit() (*model.ImageOcrTask, error) {
	if result, err := i.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.ImageOcrTask), nil
	}
}

func (i imageOcrTaskDo) FirstOrCreate() (*model.ImageOcrTask, error) {
	if result, err := i.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.ImageOcrTask), nil
	}
}

func (i imageOcrTaskDo) FindByPage(offset int, limit int) (result []*model.ImageOcrTask, count int64, err error) {
	result, err = i.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = i.Offset(-1).Limit(-1).Count()
	return
}

func (i imageOcrTaskDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = i.Count()
	if err != nil {
		return
	}

	err = i.Offset(offset).Limit(limit).Scan(result)
	return
}

func (i imageOcrTaskDo) Scan(result interface{}) (err error) {
	return i.DO.Scan(result)
}

func (i imageOcrTaskDo) Delete(models ...*model.ImageOcrTask) (result gen.ResultInfo, err error) {
	return i.DO.Delete(models)
}

func (i *imageOcrTaskDo) withDO(do gen.Dao) *imageOcrTaskDo {
	i.DO = *do.(*gen.DO)
	return i
}
