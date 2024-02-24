package _map

import "gorm.io/gorm"

type Map struct {
	gorm.Model
	EntityId    string `json:"entity_id"`
	Name        string `json:"name"`
	ReferenceId uint   `json:"reference_id"`
	EntityName  string `json:"entity_name"`
}

type RoleMap struct {
	gorm.Model
	EntityId   string `json:"entity_id" gorm:"default:null"`
	Name       string `json:"name"`
	EntityName string `json:"entity_name"`
	RoleId     uint64 `json:"role_id"`
}
