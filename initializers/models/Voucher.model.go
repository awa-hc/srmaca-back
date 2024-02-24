// models/voucher.go
package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Voucher struct {
	gorm.Model
	UserID        uint           `json:"user_id"`
	Glosa         string         `json:"glosa"`
	Img           string         `json:"img"`
	PaymentMethod string         `json:"payment_method"`
	Products      datatypes.JSON `json:"product_details"`
	TotalPrice    uint           `json:"TotalPrice"`
	Users         Users          `gorm:"foreignKey:UserID"`
	Status        bool           `json:"status"`
	Delivery      bool           `json:"delivery"`
}
