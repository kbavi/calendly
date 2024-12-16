package models

import (
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	ID    string
	Email string
	Name  string
}

func (UserModel) TableName() string {
	return "users"
}
