package models

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	FullName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Role     string `default:"user"`
}
