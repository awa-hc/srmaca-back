package handlers

import (
	"backend/initializers/database"
	"backend/initializers/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func GetUserById(c *gin.Context) {
	var user models.Users
	id := c.Param("id")
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "failed to get the user"})
		return
	}
	c.JSON(http.StatusOK, user)

}
