package handlers

import (
	"backend/initializers/database"
	"backend/initializers/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var order models.Orders
	user, _ := c.Get("user")
	UserID := user.(models.Users).ID

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{"error": "json invalid"})
		return
	}
	order.UsersID = UserID

	isDeliverySent := c.Request.PostFormValue("delivery") != ""
	if !isDeliverySent {
		c.JSON(400, gin.H{"error": "delivery not sent"})
		return
	}

	result := database.DB.Create(&order)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "failed to create the order"})
		return
	}
	c.JSON(200, gin.H{"message": "order created successfully"})
}

func GetOrders(c *gin.Context) {
	var orders []models.Orders
	if err := database.DB.Preload("Users").Find(&orders).Error; err != nil {
		c.JSON(400, gin.H{"error": "failed to get the orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}
func GetOrdersByUserID(c *gin.Context) {
	var orders []models.Orders
	user, _ := c.Get("user")
	UserID := user.(models.Users).ID
	if err := database.DB.Where("users_id = ?", UserID).Find(&orders).Error; err != nil {
		c.JSON(400, gin.H{"error": "failed to get the orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}
func GetOrder(c *gin.Context) {
	var order models.Orders
	id := c.Param("id")
	if err := database.DB.Where("id = ?", id).First(&order).Error; err != nil {
		c.JSON(400, gin.H{"error": "failed to get the order"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func DeteleOrder(c *gin.Context) {
	var order models.Orders
	id := c.Param("id")
	if err := database.DB.Where("id = ?", id).Delete(&order).Error; err != nil {
		c.JSON(400, gin.H{"error": "failed to delete the order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order deleted successfully"})
}

func UpdateOrder(c *gin.Context) {
	var order models.Orders
	id := c.Param("id")
	if err := database.DB.Where("id = ?", id).First(&order).Error; err != nil {
		c.JSON(400, gin.H{"error": "failed to get the order"})
		return
	}
	if order.Status {
		c.JSON(400, gin.H{"error": "order already completed"})
		return
	}
	if !(order.Status) {
		order.Status = true
	}
	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(400, gin.H{"error": "failed to update the order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order updated successfully"})
}
