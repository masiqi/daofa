// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dal

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

var (
	Q                = new(Query)
	Admin            *admin
	ExerciseMaterial *exerciseMaterial
	ExerciseQuestion *exerciseQuestion
	KnowledgePoint   *knowledgePoint
	Subject          *subject
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	Admin = &Q.Admin
	ExerciseMaterial = &Q.ExerciseMaterial
	ExerciseQuestion = &Q.ExerciseQuestion
	KnowledgePoint = &Q.KnowledgePoint
	Subject = &Q.Subject
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:               db,
		Admin:            newAdmin(db, opts...),
		ExerciseMaterial: newExerciseMaterial(db, opts...),
		ExerciseQuestion: newExerciseQuestion(db, opts...),
		KnowledgePoint:   newKnowledgePoint(db, opts...),
		Subject:          newSubject(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	Admin            admin
	ExerciseMaterial exerciseMaterial
	ExerciseQuestion exerciseQuestion
	KnowledgePoint   knowledgePoint
	Subject          subject
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:               db,
		Admin:            q.Admin.clone(db),
		ExerciseMaterial: q.ExerciseMaterial.clone(db),
		ExerciseQuestion: q.ExerciseQuestion.clone(db),
		KnowledgePoint:   q.KnowledgePoint.clone(db),
		Subject:          q.Subject.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:               db,
		Admin:            q.Admin.replaceDB(db),
		ExerciseMaterial: q.ExerciseMaterial.replaceDB(db),
		ExerciseQuestion: q.ExerciseQuestion.replaceDB(db),
		KnowledgePoint:   q.KnowledgePoint.replaceDB(db),
		Subject:          q.Subject.replaceDB(db),
	}
}

type queryCtx struct {
	Admin            IAdminDo
	ExerciseMaterial IExerciseMaterialDo
	ExerciseQuestion IExerciseQuestionDo
	KnowledgePoint   IKnowledgePointDo
	Subject          ISubjectDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		Admin:            q.Admin.WithContext(ctx),
		ExerciseMaterial: q.ExerciseMaterial.WithContext(ctx),
		ExerciseQuestion: q.ExerciseQuestion.WithContext(ctx),
		KnowledgePoint:   q.KnowledgePoint.WithContext(ctx),
		Subject:          q.Subject.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
