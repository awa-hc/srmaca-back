package handlers

import (
	"backend/initializers/database"
	"backend/initializers/models"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	var Request models.Product
	if err := c.ShouldBindJSON(&Request); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	if Request.Name == "" {
		c.JSON(400, gin.H{"error": "name is required"})
		return
	}
	if Request.Description == "" {
		c.JSON(400, gin.H{"error": "description is required"})
		return
	}

	if Request.Price == 0 {
		c.JSON(400, gin.H{"error": "price is required"})
		return
	}
	if Request.Stock == 0 {
		c.JSON(400, gin.H{"error": "stock is required"})
		return
	}
	if err := database.DB.Create(&Request).Error; err != nil {
		c.JSON(500, gin.H{"error": "error creating product"})
		return
	}

	c.JSON(200, gin.H{"product": Request})
}

func GetProducts(c *gin.Context) {
	var products []models.Product
	if err := database.DB.Find(&products).Error; err != nil {
		c.JSON(500, gin.H{"error": "error getting products"})
		return
	}
	c.JSON(200, gin.H{"products": products})
}

func GetProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "error getting product"})
		return
	}
	c.JSON(200, gin.H{"product": product})
}
func UpdateProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "error getting product"})
		return
	}
	var Request models.Product
	if err := c.ShouldBindJSON(&Request); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	if Request.Name == "" {
		c.JSON(400, gin.H{"error": "name is required"})
		return
	}
	if Request.Description == "" {
		c.JSON(400, gin.H{"error": "description is required"})
		return
	}

	if Request.Price == 0 {
		c.JSON(400, gin.H{"error": "price is required"})
		return
	}
	if Request.Stock == 0 {
		c.JSON(400, gin.H{"error": "stock is required"})
		return
	}
	product.Name = Request.Name
	product.Description = Request.Description
	product.Price = Request.Price
	product.Stock = Request.Stock
	product.Img = Request.Img
	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(500, gin.H{"error": "error updating product"})
		return
	}
	c.JSON(200, gin.H{"product": product})
}
