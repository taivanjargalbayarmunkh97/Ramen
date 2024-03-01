package user

import (
	"example.com/ramen/models/file"
	_map "example.com/ramen/models/map"
	"example.com/ramen/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID                 uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name               string         `json:"name" default:"null"`
	Email              string         `json:"email" default:"null" gorm:"uniqueIndex"`
	PhoneNumber        string         `json:"phone_number" default:"null"`
	Password           string         `json:"password" default:"null"`
	Provider           *string        `json:"provider" default:"null"`
	Photo              *file.File     `json:"photo" default:"null" gorm:"foreignKey:ParentId"`
	Photo1             *file.File     `json:"photo1" default:"null" gorm:"foreignKey:CompanyParentId"`
	Photo2             *file.File     `json:"photo2" default:"null" gorm:"foreignKey:InfluencerParentId;references:ID"`
	PRole              *_map.Map      `json:"prole" default:"null" gorm:"foreignKey:EntityId"`
	Role               *_map.RoleMap  `json:"role" default:"null" gorm:"foreignKey:EntityId"`
	Followers          *float64       `json:"followers" default:"null"`
	Location           *string        `json:"location" default:"null"`
	EngagementRate     *float64       `json:"engagement_rate" default:"null"`
	AverageLikes       *float64       `json:"average_likes" default:"null"`
	Bio                *string        `json:"bio" default:"null"`
	TotalPosts         *float64       `json:"total_posts" default:"null"`
	AvgLikes           *float64       `json:"avg_likes" default:"null"`
	AvgComments        *float64       `json:"avg_comments" default:"null"`
	AvgViews           *float64       `json:"avg_views" default:"null"`
	AvgReelPlays       *float64       `json:"avg_reel_plays" default:"null"`
	GenderSplit        *string        `json:"gender_split" default:"null"`
	AudienceInterests  *string        `json:"audience_interests" default:"null"`
	PopularPosts       *string        `json:"popular_posts" default:"null"`
	InfluencerIgName   *string        `json:"influencer_ig_name" default:"null"`
	CompanyAccount     *string        `json:"company_account" default:"null"`
	ManagerPhoneNumber *string        `json:"manager_phone_number" default:"null"`
	CreatedAt          *time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          *time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type UserSimple struct {
	ID                 *uuid.UUID    `json:"id"`
	Name               string        `json:"name" default:"null"`
	Email              string        `json:"email" default:"null" gorm:"uniqueIndex"`
	PhoneNumber        string        `json:"phone_number" default:"null"`
	Provider           *string       `json:"provider" default:"null"`
	Photo              *file.File    `json:"photo" default:"null" gorm:"foreignKey:ParentId"`
	Photo1             *file.File    `json:"photo1" default:"null" gorm:"foreignKey:CompanyParentId"`
	Photo2             *file.File    `json:"photo2" default:"null" gorm:"foreignKey:InfluencerParentId;references:ID"`
	PRole              *_map.Map     `json:"prole" default:"null" gorm:"foreignKey:EntityId;references:ID"`
	Role               *_map.RoleMap `json:"role" default:"null" gorm:"foreignKey:EntityId"`
	Followers          *float64      `json:"followers" default:"null"`
	Location           *string       `json:"location" default:"null"`
	EngagementRate     *float64      `json:"engagement_rate" default:"null"`
	AverageLikes       *float64      `json:"average_likes" default:"null"`
	Bio                *string       `json:"bio" default:"null"`
	TotalPosts         *float64      `json:"total_posts" default:"null"`
	AvgLikes           *float64      `json:"avg_likes" default:"null"`
	AvgComments        *float64      `json:"avg_comments" default:"null"`
	AvgViews           *float64      `json:"avg_views" default:"null"`
	AvgReelPlays       *float64      `json:"avg_reel_plays" default:"null"`
	GenderSplit        *string       `json:"gender_split" default:"null"`
	AudienceInterests  *string       `json:"audience_interests" default:"null"`
	PopularPosts       *string       `json:"popular_posts" default:"null"`
	InfluencerIgName   *string       `json:"influencer_ig_name" default:"null"`
	CompanyAccount     *string       `json:"company_account" default:"null"`
	ManagerPhoneNumber *string       `json:"manager_phone_number" default:"null"`
}

type UserUpdate struct {
	Name               string              `json:"name"`
	Email              string              `json:"email"`
	PhoneNumber        string              `json:"phone_number"`
	Role               *string             `json:"role"`
	Provider           *string             `json:"provider"`
	Photo              *utils.Base64Struct `json:"photo"`
	Followers          *float64            `json:"followers"`
	Location           *string             `json:"location"`
	EngagementRate     *float64            `json:"engagement_rate"`
	AverageLikes       *float64            `json:"average_likes"`
	Bio                *string             `json:"bio"`
	TotalPosts         *float64            `json:"total_posts"`
	AvgLikes           *float64            `json:"avg_likes"`
	AvgComments        *float64            `json:"avg_comments"`
	AvgViews           *float64            `json:"avg_views"`
	AvgReelPlays       *float64            `json:"avg_reel_plays"`
	GenderSplit        *string             `json:"gender_split"`
	AudienceInterests  *string             `json:"audience_interests"`
	PopularPosts       *string             `json:"popular_posts"`
	InfluencerIgName   *string             `json:"influencer_ig_name"`
	CompanyAccount     *string             `json:"company_account"`
	ManagerPhoneNumber *string             `json:"manager_phone_number"`
	Password           string              `json:"password" default:"null"`
}

func (u User) FilterUserRecord(user *User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		//Role:  user.Role,
		//Photo:     *user.Photo,
		CreatedAt: *user.CreatedAt,
		UpdatedAt: *user.UpdatedAt,
	}
}

type SignUpInput struct {
	Name            string             `json:"name" validate:"required"`
	Email           string             `json:"email" validate:"required"`
	Password        string             `json:"password" validate:"required,min=8"`
	PasswordConfirm string             `json:"passwordConfirm" validate:"required,min=8"`
	Photo           utils.Base64Struct `json:"photo"`
}

type SignUpInfluencer struct {
	Name            string             `json:"name" validate:"required"`
	IgName          *string            `json:"ig_name" validate:"required"`
	Email           string             `json:"email" validate:"required"`
	Password        string             `json:"password" validate:"required"`
	PasswordConfirm string             `json:"passwordConfirm" validate:"required"`
	Photo           utils.Base64Struct `json:"photo"`
	Followers       *float64           `json:"followers"`
	Location        *string            `json:"location"`
	EngagementRate  *float64           `json:"engagement_rate"`
	AverageLikes    *float64           `json:"average_likes"`
	Bio             *string            `json:"bio"`
	TotalPosts      *float64           `json:"total_posts"`
	AvgLikes        *float64           `json:"avg_likes"`
	AvgComments     *float64           `json:"avg_comments"`
	AvgViews        *float64           `json:"avg_views"`
	AvgReelPlays    *float64           `json:"avg_reel_plays"`
	GenderSplit     *string            `json:"gender_split"`
	AudienceInteres *string            `json:"audience_interests"`
	PopularPosts    *string            `json:"popular_posts"`
	RoleId          *string            `json:"role_id"`
	PhoneNumber     string             `json:"phone_number"`
}

type SignUpCompany struct {
	Name               string             `json:"name" validate:"required"`
	Email              string             `json:"email" validate:"required"`
	PhoneNumber        string             `json:"phone_number" validate:"required"`
	Password           string             `json:"password" validate:"required"`
	PasswordConfirm    string             `json:"passwordConfirm" validate:"required"`
	Photo              utils.Base64Struct `json:"photo"`
	RoleId             *string            `json:"role_id"`
	Location           *string            `json:"location"`
	CompanyAccount     *string            `json:"company_account"`
	ManagerPhoneNumber *string            `json:"manager_phone_number"`
	ProleId            string             `json:"prole_id"`
}

type SignInInput struct {
	Email    string `json:"email"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      *string   `json:"role,omitempty"`
	Photo     file.File `json:"photo,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
