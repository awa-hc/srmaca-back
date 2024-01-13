package models

import "gorm.io/gorm"

type Orders struct {
	gorm.Model
	OrderID    string `json:"order_id"`
	OrderDate  string `json:"order_date"`
	OrderTotal string `json:"order_total"`
	Delivery   bool   `json:"delivery"`
	UsersID    uint   `json:"users_id"`
	Users      Users  `gorm:"foreignKey:UsersID"`
	Status     bool   `json:"status"`
}
