package reference

import "gorm.io/gorm"

type Reference struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Field1      string `json:"field1"`
	Field2      string `json:"field2"`
	Field3      string `json:"field3"`
	Code        string `json:"code"`
}

type CreateReference struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Field1      string `json:"field1"`
	Field2      string `json:"field2"`
	Field3      string `json:"field3"`
	Code        string `json:"code" validate:"required"`
}

type UpdateReference struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Field1      string `json:"field1"`
	Field2      string `json:"field2"`
	Field3      string `json:"field3"`
	Code        string `json:"code" validate:"required"`
}
