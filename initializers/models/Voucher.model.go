// models/voucher.go
package models

import "gorm.io/gorm"

type Voucher struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	ProductID uint   `json:"product_id"`
	Glosa     string `json:"glosa"`
	Img       string `json:"img"`
	Users     Users  `gorm:"foreignKey:UserID"`
}
