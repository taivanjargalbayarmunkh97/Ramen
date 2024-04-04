package news

import (
	"example.com/ramen/models/file"
	_map "example.com/ramen/models/map"
	"example.com/ramen/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type News struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Body        string         `json:"body"`
	Image       file.File      `json:"image" gorm:"foreignKey:NewsParentId;references:ID"`
	Type        []_map.Map     `json:"type" gorm:"foreignKey:NewsEntityId;references:ID"`
}

type CreateNews struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Image       utils.Base64Struct
	Body        string   `json:"body"`
	Type        []string `json:"type"`
}

type UpdateNews struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Image       utils.Base64Struct
	Body        string   `json:"body"`
	Type        []string `json:"type"`
}
