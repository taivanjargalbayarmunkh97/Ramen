package channel

import (
	"example.com/ramen/models/file"
	_map "example.com/ramen/models/map"
	"example.com/ramen/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Channel struct {
	ID                uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt         time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	Email             string         `json:"email"`
	Phone             string         `json:"phone"`
	Website           string         `json:"website"`
	Address           string         `json:"address"`
	TvDailyAvgViews   string         `json:"tv_daily_avg_views"`
	TvUnivisionNumber string         `json:"tv_univision_number"`
	FmDailyAvg1       string         `json:"fm_daily_avg_1"`
	FmDailyAvg2       string         `json:"fm_daily_avg_2"`
	FmSecondEval1     string         `json:"fm_second_eval_1"`
	FmSecondEval2     string         `json:"fm_second_eval_2"`
	FmCpe1            string         `json:"fm_cpe_1"`
	FmCpe2            string         `json:"fm_cpe_2"`
	Image             file.File      `json:"image" gorm:"foreignKey:ChannelParentId;references:ID"`
	Type              _map.Map       `json:"type" gorm:"foreignKey:ChannelTypeEntityId;references:ID"`
}

type CreateChannel struct {
	Type              string             `json:"type" validate:"required"`
	Name              string             `json:"name" validate:"required"`
	Description       string             `json:"description"`
	Email             string             `json:"email"`
	Phone             string             `json:"phone"`
	Website           string             `json:"website"`
	Address           string             `json:"address"`
	Image             utils.Base64Struct `json:"image" validate:"required"`
	TvDailyAvgViews   string             `json:"tv_daily_avg_views"`
	TvUnivisionNumber string             `json:"tv_univision_number"`
	FmDailyAvg1       string             `json:"fm_daily_avg_1"`
	FmDailyAvg2       string             `json:"fm_daily_avg_2"`
	FmSecondEval1     string             `json:"fm_second_eval_1"`
	FmSecondEval2     string             `json:"fm_second_eval_2"`
	FmCpe1            string             `json:"fm_cpe_1"`
	FmCpe2            string             `json:"fm_cpe_2"`
}

type UpdateChannel struct {
	Type              string             `json:"type" validate:"required"`
	Name              string             `json:"name" validate:"required"`
	Description       string             `json:"description"`
	Email             string             `json:"email"`
	Phone             string             `json:"phone"`
	Website           string             `json:"website"`
	Address           string             `json:"address"`
	Image             utils.Base64Struct `json:"image" validate:"required"`
	TvDailyAvgViews   string             `json:"tv_daily_avg_views"`
	TvUnivisionNumber string             `json:"tv_univision_number"`
	FmDailyAvg1       string             `json:"fm_daily_avg_1"`
	FmDailyAvg2       string             `json:"fm_daily_avg_2"`
	FmSecondEval1     string             `json:"fm_second_eval_1"`
	FmSecondEval2     string             `json:"fm_second_eval_2"`
	FmCpe1            string             `json:"fm_cpe_1"`
	FmCpe2            string             `json:"fm_cpe_2"`
}
