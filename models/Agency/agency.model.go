package Agency

import (
	"example.com/ramen/models/file"
	_map "example.com/ramen/models/map"
	"example.com/ramen/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Agency struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Body        string         `json:"body"`
	Website     string         `json:"website"`
	Email       string         `json:"email"`
	Phone       string         `json:"phone"`
	Address     string         `json:"address"`
	City        string         `json:"city"`
	Type        []_map.Map     `json:"type" gorm:"foreignKey:AgencyEntityId"`
	Brands      []_map.Map     `json:"brands" gorm:"foreignKey:AgencyBrandsEntityId"`
	Image       file.File      `json:"image" gorm:"foreignKey:AgencyParentId"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type CreateAgency struct {
	Name        string             `json:"name" validate:"required"`
	Description string             `json:"description" validate:"required"`
	Body        string             `json:"body" validate:"required"`
	Website     string             `json:"website" validate:"required"`
	Email       string             `json:"email" validate:"required"`
	Phone       string             `json:"phone" validate:"required"`
	Address     string             `json:"address" validate:"required"`
	City        string             `json:"city" validate:"required"`
	Image       utils.Base64Struct `json:"image" validate:"required"`
	Type        []string           `json:"type" validate:"required"`
	Brands      []string           `json:"brands" validate:"required"`
}

type UpdateAgency struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Body        string             `json:"body"`
	Website     string             `json:"website"`
	Email       string             `json:"email"`
	Phone       string             `json:"phone"`
	Address     string             `json:"address"`
	City        string             `json:"city"`
	Image       utils.Base64Struct `json:"image"`
	Type        []string           `json:"type"`
	Brands      []string           `json:"brands"`
}
