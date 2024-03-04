package resources

import (
	"example.com/ramen/models/file"
	_map "example.com/ramen/models/map"
	"example.com/ramen/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Resources struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Body          string         `json:"body"`
	YoutubeLink   string         `json:"youtube_link"`
	FacebookLink  string         `json:"facebook_link"`
	TwitterLink   string         `json:"twitter_link"`
	InstagramLink string         `json:"instagram_link"`
	LinkedinLink  string         `json:"linkedin_link"`
	PinterestLink string         `json:"pinterest_link"`
	Image         file.File      `json:"image" gorm:"foreignKey:ResourceParentId;references:ID"`
	Type          []_map.Map     `json:"type" gorm:"foreignKey:ResourcesEntityId;references:ID"`
}

type CreateResources struct {
	Name          string             `json:"name" validate:"required,max=200"`
	Description   string             `json:"description,max=1000"`
	Body          string             `json:"body"`
	YoutubeLink   string             `json:"youtube_link"`
	FacebookLink  string             `json:"facebook_link"`
	TwitterLink   string             `json:"twitter_link"`
	InstagramLink string             `json:"instagram_link"`
	LinkedinLink  string             `json:"linkedin_link"`
	PinterestLink string             `json:"pinterest_link"`
	Image         utils.Base64Struct `json:"image"`
	Type          []string           `json:"type"`
}

type UpdateResources struct {
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	Body          string             `json:"body"`
	YoutubeLink   string             `json:"youtube_link"`
	FacebookLink  string             `json:"facebook_link"`
	TwitterLink   string             `json:"twitter_link"`
	InstagramLink string             `json:"instagram_link"`
	LinkedinLink  string             `json:"linkedin_link"`
	PinterestLink string             `json:"pinterest_link"`
	Image         utils.Base64Struct `json:"image"`
	Type          []string           `json:"type"`
}
