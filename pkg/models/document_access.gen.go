// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

import (
	"time"
)

const TableNameDocumentAccess = "document_access"

// DocumentAccess mapped from table <document_access>
type DocumentAccess struct {
	ID             int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	DocumentID     string    `gorm:"column:document_id;not null" json:"document_id"`
	UserID         string    `gorm:"column:user_id;not null" json:"user_id"`
	AccessLevel    string    `gorm:"column:access_level;not null" json:"access_level"`
	LastAccessedAt time.Time `gorm:"column:last_accessed_at" json:"last_accessed_at"`
}

// TableName DocumentAccess's table name
func (*DocumentAccess) TableName() string {
	return TableNameDocumentAccess
}
