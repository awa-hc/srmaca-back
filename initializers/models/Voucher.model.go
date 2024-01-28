// models/voucher.go
package models

import "gorm.io/gorm"

type Voucher struct {
	gorm.Model
	UserID        uint   `json:"user_id"`
	ProductID     uint   `json:"product_id"`
	Glosa         string `json:"glosa"`
	Img           string `json:"img"`
	PaymentMethod string `json:"payment_method"`
	Quantity      uint   `json:"quantity"`
	Amount        uint   `json:"amount"`
	Users         Users  `gorm:"foreignKey:UserID"`
	Status        bool   `json:"status"`
}
