// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameKnowledgePoint = "knowledge_point"

// KnowledgePoint mapped from table <knowledge_point>
type KnowledgePoint struct {
	ID          int32      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	SubjectID   int32      `gorm:"column:subject_id;not null" json:"subject_id"`
	ParentID    *int32     `gorm:"column:parent_id" json:"parent_id"`
	Name        string     `gorm:"column:name;not null" json:"name"`
	Description *string    `gorm:"column:description" json:"description"`
	IsLeaf      bool       `gorm:"column:is_leaf;not null" json:"is_leaf"`
	CreatedAt   *time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName KnowledgePoint's table name
func (*KnowledgePoint) TableName() string {
	return TableNameKnowledgePoint
}
