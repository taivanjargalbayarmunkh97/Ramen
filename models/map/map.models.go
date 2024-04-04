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
	ResourcesEntityId       string `json:"resources_entity_id" gorm:"default:null"`
	AgencyEntityId          string `json:"agency_entity_id" gorm:"default:null"`
	AgencyBrandsEntityId    string `json:"agency_brands_entity_id" gorm:"default:null"`
	NewsEntityId            string `json:"news_entity_id" gorm:"default:null"`
	Name                    string `json:"name"`
	ReferenceId             uint   `json:"reference_id"`
	EntityName              string `json:"entity_name"`
	Image                   string `json:"image"`
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
