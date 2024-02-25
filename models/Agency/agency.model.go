package Agency

import (
	_map "example.com/ramen/models/map"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Agency struct {
	ID          uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Website     string           `json:"website"`
	Email       string           `json:"email"`
	Phone       string           `json:"phone"`
	Address     string           `json:"address"`
	City        string           `json:"city"`
	Type        []_map.AgencyMap `json:"type" gorm:"foreignKey:EntityId"`
	CreatedAt   time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt   `json:"deleted_at" gorm:"index"`
}

type CreateAgency struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Website     string   `json:"website" validate:"required"`
	Email       string   `json:"email" validate:"required"`
	Phone       string   `json:"phone" validate:"required"`
	Address     string   `json:"address" validate:"required"`
	City        string   `json:"city" validate:"required"`
	Type        []string `json:"type" validate:"required"`
}

type UpdateAgency struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	City        string `json:"city"`
}
