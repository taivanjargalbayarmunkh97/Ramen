package Company

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Company struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Website     string         `json:"website"`
	Email       string         `json:"email"`
	Phone       string         `json:"phone"`
	Address     string         `json:"address"`
	//Image           file.File      `json:"image" gorm:"foreignKey:CompanyParentId;references:ID"`
	City            string `json:"city"`
	AreasOfActivity string `json:"areas_of_activity"`
}

type CreateCompany struct {
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description" validate:"required"`
	Website         string `json:"website" validate:"required"`
	Email           string `json:"email" validate:"required"`
	Phone           string `json:"phone" validate:"required"`
	Address         string `json:"address" validate:"required"`
	City            string `json:"city" validate:"required"`
	Image           string `json:"image" validate:"required"`
	AreasOfActivity string `json:"areas_of_activity" validate:"required"`
}

type UpdateCompany struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Website         string `json:"website"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Address         string `json:"address"`
	City            string `json:"city"`
	Image           string `json:"image"`
	AreasOfActivity string `json:"areas_of_activity"`
}
