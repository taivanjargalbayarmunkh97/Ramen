package file

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	ID                 uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;autoUpdateTime:false"`
	DeletedAt          gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	ParentId           string         `json:"parent_id" gorm:"default:null"`
	CompanyParentId    string         `json:"company_parent_id" gorm:"default:null"`
	InfluencerParentId string         `json:"influencer_parent_id" gorm:"default:null"`
	ChannelParentId    string         `json:"channel_parent_id" gorm:"default:null"`
	FileName           string         `json:"file_name"`
	FilePath           string         `json:"file_path"`
	Size               string         `json:"size"`
	Category           string         `json:"category"`
}
