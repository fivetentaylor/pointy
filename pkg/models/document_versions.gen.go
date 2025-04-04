// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

import (
	"time"
)

const TableNameDocumentVersion = "document_versions"

// DocumentVersion mapped from table <document_versions>
type DocumentVersion struct {
	ID             string    `gorm:"column:id;primaryKey;default:uuid_generate_v4()" json:"id"`
	DocumentID     string    `gorm:"column:document_id;not null" json:"document_id"`
	Name           string    `gorm:"column:name;not null" json:"name"`
	ContentAddress string    `gorm:"column:content_address;not null" json:"content_address"`
	CreatedAt      time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy      string    `gorm:"column:created_by;not null" json:"created_by"`
	UpdatedBy      string    `gorm:"column:updated_by;not null" json:"updated_by"`
}

// TableName DocumentVersion's table name
func (*DocumentVersion) TableName() string {
	return TableNameDocumentVersion
}
