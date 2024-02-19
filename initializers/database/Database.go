package database

import (
	"backend/initializers/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	cs := os.Getenv("DB")
	DB, err = gorm.Open(postgres.Open(cs), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB.AutoMigrate(&models.Users{}, &models.Voucher{}, &models.Product{})

}
