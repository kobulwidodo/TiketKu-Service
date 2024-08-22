package entity

import (
	"go-clean/src/lib/auth"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string
	Password string `json:"-"`
	Name     string
	IsAdmin  bool
}

type UserParam struct {
	ID    uint
	Email string
}

type CreateUserParam struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
	Name     string `binding:"required"`
}

type LoginUserParam struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) ConvertToAuthUser() auth.User {
	return auth.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
		IsAdmin:  u.IsAdmin,
	}
}
