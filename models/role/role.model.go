package role

import (
	"gorm.io/gorm"
	"time"
)

type Role struct {
	Id          uint64         `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key;uniqueIndex"`
	Name        string         `json:"name" gorm:"default:null"`
	Description string         `json:"description"`
	Field1      string         `json:"field1"`
	Field2      string         `json:"field2"`
	Field3      string         `json:"field3"`
	CreatedAt   *time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   *time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type RoleCreateInput struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Field1      string `json:"field1" default:"null"`
	Field2      string `json:"field2" default:"null"`
	Field3      string `json:"field3" default:"null"`
}

type RoleUpdateInput struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Field1      string `json:"field1"`
	Field2      string `json:"field2"`
	Field3      string `json:"field3"`
}
