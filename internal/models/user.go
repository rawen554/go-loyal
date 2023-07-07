package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        uint64    `gorm:"primaryKey" json:"id,omitempty"`
	Login     string    `gorm:"varchar(100);index:idx_name,unique" json:"login"`
	Password  string    `gorm:"type:varchar(255);not null"`
	LastLogin time.Time `gorm:"default:null;" json:"lastLogin,omitempty"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time `gorm:"not null;default:now();ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt"`
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

type UserCredentialsSchema struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserSchema struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}
