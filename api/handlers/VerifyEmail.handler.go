// handlers/handlers.go

package handlers

import (
	"backend/api/email"
	"backend/initializers/database"
	"backend/initializers/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyEmail(c *gin.Context) {
	token := c.Param("token")

	// Verificar el token
	userID, err := email.VerifyVerificationToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification token"})
		return
	}

	// Marcar el usuario como verificado en la base de datos
	var user models.Users
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user"})
		return
	}

	// Verificar si el correo ya ha sido verificado
	if user.EmailVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already verified"})
		return
	}

	// Marcar el correo como verificado
	user.EmailVerified = true
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}
