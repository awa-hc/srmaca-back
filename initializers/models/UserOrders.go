package models

type UserOrders struct {
	OrderID []int `json:"order_id"`
	UserID  uint  `json:"user_id"`
	User    Users `gorm:"foreignKey:UserID"`
}
