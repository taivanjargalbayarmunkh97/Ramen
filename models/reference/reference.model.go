package reference

import (
	"example.com/ramen/models/file"
	"example.com/ramen/utils"
	"gorm.io/gorm"
)

type Reference struct {
	gorm.Model
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Field1      string     `json:"field1"`
	Field2      string     `json:"field2"`
	Field3      string     `json:"field3"`
	Code        string     `json:"code"`
	Image       *file.File `json:"image" gorm:"foreignKey:ReferenceParentId"`
}

type CreateReference struct {
	Name        string             `json:"name" validate:"required"`
	Description string             `json:"description"`
	Field1      string             `json:"field1"`
	Field2      string             `json:"field2"`
	Field3      string             `json:"field3"`
	Code        string             `json:"code" validate:"required"`
	Image       utils.Base64Struct `json:"image"`
}

type UpdateReference struct {
	Name        string             `json:"name" validate:"required"`
	Description string             `json:"description"`
	Field1      string             `json:"field1"`
	Field2      string             `json:"field2"`
	Field3      string             `json:"field3"`
	Code        string             `json:"code" validate:"required"`
	Image       utils.Base64Struct `json:"image"`
}
