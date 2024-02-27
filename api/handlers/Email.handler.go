// handlers/handlers.go

package handlers

import (
	"backend/api/email"
	"backend/initializers/database"
	"backend/initializers/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContactForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func ContactEmail(c *gin.Context) {
	var form ContactForm

	if err := c.BindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Falla en parsear el cuerpo de la solicitud❌"})
		return
	}

	// Enviar el correo
	if err := email.SendMailContact(form.Name, form.Email, form.Phone, form.Subject, form.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falla en enviar el correo❌"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Correo Enviado Exitosamente✔️"})
}

func VerifyEmail(c *gin.Context) {
	token := c.Param("token")

	// Verificar el token
	userID, err := email.VerifyVerificationToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token De Verificación Inválido❌"})
		return
	}

	// Marcar el usuario como verificado en la base de datos
	var user models.Users
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falla en encontrar el usuario❌"})
		return
	}

	// Verificar si el correo ya ha sido verificado
	if user.EmailVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "❕El correo ya se encuentra verificado❕, no se requiere verificación👍. Por favor inicie sesión.⤴️"})
		return
	}

	// Marcar el correo como verificado
	user.EmailVerified = true
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falla en verificar el usuario❌"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Correo Verificado Exitosamente ✔️ Por Favor Inicie Sesión"})
}
