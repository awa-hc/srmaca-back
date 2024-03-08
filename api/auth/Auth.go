// auth.go

package auth

import (
	"backend/api/email"
	"backend/initializers/database"
	"backend/initializers/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var body struct {
		Fullname        string `json:"fullname"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		Phone           string `json:"phone"`
		ConfirmPassword string `json:"confirmpassword"`
		Address         string `json:"address"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json invalid"})
		return
	}
	if len(body.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 8 characters"})
		return
	}

	if body.ConfirmPassword != body.Password {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password do not match"})
		return
	}
	if body.Fullname == "" || body.Email == "" || body.Password == "" || body.Address == "" || body.ConfirmPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty fields"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to hash the password"})
		return
	}

	user := models.Users{
		Email:         body.Email,
		Password:      string(hash),
		FullName:      body.Fullname,
		Phone:         body.Phone,
		Address:       body.Address,
		EmailVerified: false}

	user.Role = "user"

	result := database.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create the user"})
		return
	}

	// Generar token de verificación
	verificationToken, err := email.GenerateVerificationToken(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate verification token"})
		return
	}

	go email.SendVerificationEmail(user.Email, user.FullName, verificationToken)

	c.JSON(http.StatusOK, gin.H{"message": "user created successfully, check your email for verification", "user": user})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json invalid"})
		return
	}

	var user models.Users
	if err := database.DB.First(&user, "email = ?", body.Email).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create token"})
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Auth", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	// id := user.(models.Users).ID
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func ValidateAdmin(C *gin.Context) {
	user, _ := C.Get("user")
	if user.(models.Users).Role != "admin" {
		C.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	C.JSON(http.StatusOK, gin.H{"user": user})
}

type ResetPasswordRequest struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func ResetPassword(c *gin.Context) {
	var body ResetPasswordRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json invalid"})
		return
	}

	// Verificar si la contraseña y la confirmación coinciden
	if body.Password != body.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	// Obtener el ID de usuario del contexto
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
		return
	}
	uid := userID.(uint)

	// Generar el hash de la nueva contraseña
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Actualizar la contraseña del usuario en la base de datos
	if err := database.DB.Model(&models.Users{}).Where("id = ?", uid).Update("password", string(hash)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

func ForgotPassword(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json invalid"})
		return
	}

	var user models.Users
	if err := database.DB.First(&user, "email = ?", body.Email).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	// Generar token de restablecimiento de contraseña
	resetToken, err := email.GenerateResetPasswordToken(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate reset password token"})
		return
	}

	// Envía el correo electrónico con el token de restablecimiento de contraseña
	go email.SendResetPasswordEmail(user.Email, resetToken)

	c.JSON(http.StatusOK, gin.H{"message": "reset password email sent"})
}
