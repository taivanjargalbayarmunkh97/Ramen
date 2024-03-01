package _map

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Map struct {
	gorm.Model
	EntityId                string `json:"entity_id" gorm:"default:null"`
	CompanyActivityEntityId string `json:"company_activity_entity_id" gorm:"default:null"`
	ChannelTypeEntityId     string `json:"channel_type_entity_id" gorm:"default:null"`
	Name                    string `json:"name"`
	ReferenceId             uint   `json:"reference_id"`
	EntityName              string `json:"entity_name"`
}

type RoleMap struct {
	gorm.Model
	EntityId   string `json:"entity_id" gorm:"default:null"`
	Name       string `json:"name"`
	EntityName string `json:"entity_name"`
	RoleId     uint64 `json:"role_id"`
}

type AgencyMap struct {
	gorm.Model
	EntityId    uuid.UUID `json:"entity_id" gorm:"type:uuid"`
	Name        string    `json:"name"`
	EntityName  string    `json:"entity_name"`
	ReferenceId string    `json:"reference_id"`
}
