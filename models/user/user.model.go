package user

import (
	"gorm.io/gorm"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	ID        *uuid.UUID     `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null"`
	Email     string         `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string         `json:"password" gorm:"type:varchar(100);not null"`
	Role      *string        `json:"role" gorm:"type:varchar(50);default:'user';not null"`
	Provider  *string        `json:"provider" gorm:"type:varchar(50);default:'local';not null"`
	Photo     *string        `json:"photo" gorm:"not null;default:'default.png'"`
	Verified  *bool          `json:"verified" gorm:"not null;default:false"`
	CreatedAt *time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (u User) FilterUserRecord(user *User) UserResponse {
	return UserResponse{
		ID:        *user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      *user.Role,
		Photo:     *user.Photo,
		Provider:  *user.Provider,
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

type SignInInput struct {
	Email    string `json:"email"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	Photo     string    `json:"photo,omitempty"`
	Provider  string    `json:"provider"`
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

func FilterUserRecord(user *User) UserResponse {
	return UserResponse{
		ID:        *user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      *user.Role,
		Photo:     *user.Photo,
		Provider:  *user.Provider,
		CreatedAt: *user.CreatedAt,
		UpdatedAt: *user.UpdatedAt,
	}
}
