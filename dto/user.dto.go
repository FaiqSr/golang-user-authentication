package dto

import (
	"github.com/go-playground/validator/v10"
)

type (
	UserCreateUserRequest struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	UserLoginRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	UserUpdateRequest struct {
		Name     string `json:"name,omitempty"`
		Email    string `json:"email,omitempty"  validate:"omitempty,email"`
		Password string `json:"password,omitempty"  validate:"omitempty,min=8"`
	}
	UserResponse struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func (ucr UserCreateUserRequest) ValidateUserCreateRequest() error {
	err := validate.Struct(ucr)
	if err != nil {
		return err
	}

	return nil

}

func (ul UserLoginRequest) ValidateUserLoginRequest() error {
	err := validate.Struct(ul)
	if err != nil {
		return err
	}

	return err
}

func (uur UserUpdateRequest) ValidateUserUpdateRequest() error {
	err := validate.Struct(uur)
	if err != nil {
		return err
	}

	return nil
}
