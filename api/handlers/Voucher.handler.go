package handlers

import (
	"backend/api/utils"
	"backend/initializers/database"
	"backend/initializers/models"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateVoucher(c *gin.Context) {

	var request models.Voucher
	user, _ := c.Get("user")
	userID := user.(models.Users).ID

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Extraer el tipo de imagen y los datos base64 del campo "img"
	parts := strings.Split(request.Img, ";base64,")
	if len(parts) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format"})
		return
	}

	imgType := strings.TrimPrefix(parts[0], "data:image/")
	imgType = strings.TrimSuffix(imgType, ";base64")

	// Decodificar la imagen base64
	decodedImg, err := utils.DecodeBase64Image(parts[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode the image"})
		return
	}

	// Crear un directorio "uploads" si no existe
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create the uploads directory"})
		return
	}

	// Guardar la imagen en el directorio "uploads" con un nombre único
	filename := fmt.Sprintf("%d_%s.%s", userID, request.Glosa, imgType)
	path := filepath.Join("uploads", filename)

	if err := os.WriteFile(path, decodedImg, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
		return
	}

	voucher := models.Voucher{
		UserID:    userID,
		ProductID: request.ProductID,
		Glosa:     request.Glosa,
		Img:       filename,
	}

	if err := database.DB.Create(&voucher).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create the voucher record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Voucher uploaded successfully",
	})
}

func DeleteVoucher(c *gin.Context) {
	var voucher models.Voucher
	id := c.Param("id")
	if err := database.DB.Where("id = ?", id).Delete(&voucher).Error; err != nil {
		c.JSON(400, gin.H{"error": "failed to get the order"})
		return
	}
	c.JSON(http.StatusOK, voucher)
}
func GetVoucher(c *gin.Context) {
	var voucher models.Voucher
	id := c.Param("id")
	if err := database.DB.Where("id = ?", id).Preload("Users").First(&voucher).Error; err != nil {
		c.JSON(400, gin.H{"error": "failed to get the order"})
		return
	}
	c.JSON(http.StatusOK, voucher)
}

func GetVouchers(C *gin.Context) {
	var vouchers []models.Voucher
	if err := database.DB.Preload("Users").Find(&vouchers).Error; err != nil {
		C.JSON(400, gin.H{"error": "failed to get the vouchers"})
		return
	}
	C.JSON(http.StatusOK, vouchers)
}
