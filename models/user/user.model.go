package user

import (
	"example.com/ramen/models/file"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID                 *uuid.UUID     `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name               string         `json:"name" default:"null"`
	Email              string         `json:"email" default:"null" gorm:"uniqueIndex"`
	PhoneNumber        string         `json:"phone_number" default:"null"`
	Password           string         `json:"password" default:"null"`
	Role               *string        `json:"role" default:"null"`
	Provider           *string        `json:"provider" default:"null"`
	Photo              *[]file.File   `json:"photo" default:"null" gorm:"foreignKey:ParentId"`
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
	Name               string       `json:"name" default:"null"`
	Email              string       `json:"email" default:"null" gorm:"uniqueIndex"`
	PhoneNumber        string       `json:"phone_number" default:"null"`
	Role               *string      `json:"role" default:"null"`
	Provider           *string      `json:"provider" default:"null"`
	Photo              *[]file.File `json:"photo" default:"null" gorm:"foreignKey:ParentId"`
	Followers          *float64     `json:"followers" default:"null"`
	Location           *string      `json:"location" default:"null"`
	EngagementRate     *float64     `json:"engagement_rate" default:"null"`
	AverageLikes       *float64     `json:"average_likes" default:"null"`
	Bio                *string      `json:"bio" default:"null"`
	TotalPosts         *float64     `json:"total_posts" default:"null"`
	AvgLikes           *float64     `json:"avg_likes" default:"null"`
	AvgComments        *float64     `json:"avg_comments" default:"null"`
	AvgViews           *float64     `json:"avg_views" default:"null"`
	AvgReelPlays       *float64     `json:"avg_reel_plays" default:"null"`
	GenderSplit        *string      `json:"gender_split" default:"null"`
	AudienceInterests  *string      `json:"audience_interests" default:"null"`
	PopularPosts       *string      `json:"popular_posts" default:"null"`
	InfluencerIgName   *string      `json:"influencer_ig_name" default:"null"`
	CompanyAccount     *string      `json:"company_account" default:"null"`
	ManagerPhoneNumber *string      `json:"manager_phone_number" default:"null"`
}

type UserUpdate struct {
	Name               string   `json:"name"`
	Email              string   `json:"email"`
	PhoneNumber        string   `json:"phone_number"`
	Role               *string  `json:"role"`
	Provider           *string  `json:"provider"`
	Photo              *string  `json:"photo" `
	Followers          *float64 `json:"followers"`
	Location           *string  `json:"location"`
	EngagementRate     *float64 `json:"engagement_rate"`
	AverageLikes       *float64 `json:"average_likes"`
	Bio                *string  `json:"bio"`
	TotalPosts         *float64 `json:"total_posts"`
	AvgLikes           *float64 `json:"avg_likes"`
	AvgComments        *float64 `json:"avg_comments"`
	AvgViews           *float64 `json:"avg_views"`
	AvgReelPlays       *float64 `json:"avg_reel_plays"`
	GenderSplit        *string  `json:"gender_split"`
	AudienceInterests  *string  `json:"audience_interests"`
	PopularPosts       *string  `json:"popular_posts"`
	InfluencerIgName   *string  `json:"influencer_ig_name"`
	CompanyAccount     *string  `json:"company_account"`
	ManagerPhoneNumber *string  `json:"manager_phone_number"`
}

func (u User) FilterUserRecord(user *User) UserResponse {
	return UserResponse{
		ID:    *user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
		//Photo:     *user.Photo,
		CreatedAt: *user.CreatedAt,
		UpdatedAt: *user.UpdatedAt,
	}
}

type SignUpInput struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,min=8"`
	Photo           string `json:"photo"`
}

type SignUpInfluencer struct {
	Name            string   `json:"name" validate:"required"`
	IgName          *string  `json:"ig_name" validate:"required"`
	Email           string   `json:"email" validate:"required"`
	Password        string   `json:"password" validate:"required"`
	PasswordConfirm string   `json:"passwordConfirm" validate:"required"`
	Photo           string   `json:"photo"`
	Followers       *float64 `json:"followers"`
	Location        *string  `json:"location"`
	EngagementRate  *float64 `json:"engagement_rate"`
	AverageLikes    *float64 `json:"average_likes"`
	Bio             *string  `json:"bio"`
	TotalPosts      *float64 `json:"total_posts"`
	AvgLikes        *float64 `json:"avg_likes"`
	AvgComments     *float64 `json:"avg_comments"`
	AvgViews        *float64 `json:"avg_views"`
	AvgReelPlays    *float64 `json:"avg_reel_plays"`
	GenderSplit     *string  `json:"gender_split"`
	AudienceInteres *string  `json:"audience_interests"`
	PopularPosts    *string  `json:"popular_posts"`
	RoleId          *string  `json:"role_id"`
	PhoneNumber     string   `json:"phone_number"`
}

type SignUpCompany struct {
	Name               string  `json:"name" validate:"required"`
	Email              string  `json:"email" validate:"required"`
	PhoneNumber        string  `json:"phone_number" validate:"required"`
	Password           string  `json:"password" validate:"required"`
	PasswordConfirm    string  `json:"passwordConfirm" validate:"required"`
	Photo              string  `json:"photo"`
	RoleId             *string `json:"role_id"`
	Location           *string `json:"location"`
	CompanyAccount     *string `json:"company_account"`
	ManagerPhoneNumber *string `json:"manager_phone_number"`
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

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
