package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"math"

	"gorm.io/gorm"
)

func Paginate(value interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {

	var totalRows int64
	if err := db.Count(&totalRows).Error; err != nil {
		println(err.Error())
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

	pageSize := pagination.PerPage
	if pageSize <= 0 {
		pagination.PerPage = int(totalRows)
		pagination.TotalPages = 1
		pagination.CurrentPageNo = 1
		pagination.TotalElements = totalRows
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

	page := pagination.CurrentPageNo
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * pageSize
	pagination.PerPage = pageSize
	pagination.CurrentPageNo = page
	pagination.TotalElements = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pageSize)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Order(pagination.Sort).Limit(pageSize)
	}
}

func PaginateString(pagination *Pagination) string {
	pageSize := pagination.PerPage
	if pageSize <= 0 {
		return ""
	}
	page := pagination.CurrentPageNo
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize
	if pagination.Sort != "" {
		return fmt.Sprintf("order by %s offset %d limit %d", pagination.Sort, offset, pageSize)
	} else {
		return fmt.Sprintf("offset %d limit %d", offset, pageSize)
	}
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}
