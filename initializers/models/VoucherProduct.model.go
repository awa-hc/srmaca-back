package models

type VoucherProducts struct {
	VoucherID uint `gorm:"primaryKey"`
	ProductID uint `gorm:"primaryKey"`
	Quantity  uint `json:"quantity"`
}
