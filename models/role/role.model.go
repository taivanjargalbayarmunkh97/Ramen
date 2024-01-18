package role

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Role struct {
	ID          *uuid.UUID     `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name        string         `json:"name" gorm:"default:nullcd"`
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
	Field1      string `json:"field1"`
	Field2      string `json:"field2"`
	Field3      string `json:"field3"`
}

type RoleUpdateInput struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Field1      string `json:"field1"`
	Field2      string `json:"field2"`
	Field3      string `json:"field3"`
}
