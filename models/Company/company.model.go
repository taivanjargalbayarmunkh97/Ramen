package Company

import (
	"example.com/ramen/models/file"
	_map "example.com/ramen/models/map"
	"example.com/ramen/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Company struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	Body            string         `json:"body"`
	YoutubeLink     string         `json:"youtube_link"`
	Website         string         `json:"website"`
	Email           string         `json:"email"`
	Phone           string         `json:"phone"`
	Address         string         `json:"address"`
	Image           file.File      `json:"image" gorm:"foreignKey:CampaignsParentId;references:ID"`
	AreasOfActivity []_map.Map     `json:"areas_of_activity" gorm:"foreignKey:CompanyActivityEntityId;references:ID"`
	City            string         `json:"city"`
}

type CreateCompany struct {
	Name            string             `json:"name" validate:"required"`
	Description     string             `json:"description"`
	Body            string             `json:"body"`
	YoutubeLink     string             `json:"youtube_link"`
	Website         string             `json:"website"`
	Email           string             `json:"email" validate:"required"`
	Phone           string             `json:"phone"`
	Address         string             `json:"address"`
	City            string             `json:"city"`
	Image           utils.Base64Struct `json:"image" validate:"required"`
	AreasOfActivity []string           `json:"areas_of_activity"`
}

type UpdateCompany struct {
	Id              string             `json:"id"`
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	Body            string             `json:"body"`
	YoutubeLink     string             `json:"youtube_link"`
	Website         string             `json:"website"`
	Email           string             `json:"email"`
	Phone           string             `json:"phone"`
	Address         string             `json:"address"`
	City            string             `json:"city"`
	Image           utils.Base64Struct `json:"image"`
	AreasOfActivity []string           `json:"areas_of_activity"`
}
