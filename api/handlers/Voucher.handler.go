package handlers

import (
	"backend/api/utils"
	"backend/initializers/database"
	"backend/initializers/models"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateVoucher(c *gin.Context) {

	// Leer el cuerpo de la solicitud
	body, _ := io.ReadAll(c.Request.Body)

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var request models.Voucher
	user, _ := c.Get("user")
	userID := user.(models.Users).ID

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Si el método de pago es "Transferencia Bancaria", entonces la imagen es requerida
	if request.PaymentMethod == "transfer" {
		if request.Img == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required for bank transfer"})
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

		request.Img = filename
	} else if request.PaymentMethod == "cash" {
		request.Img = "" // Asegúrate de que la imagen esté vacía
	}

	voucher := models.Voucher{
		UserID:        userID,
		Glosa:         request.Glosa,
		PaymentMethod: request.PaymentMethod,
		Img:           request.Img,
		TotalPrice:    request.TotalPrice,
		Products:      request.Products,
		Status:        false,
	}

	if err := database.DB.Create(&voucher).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create the voucher record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Voucher uploaded successfully", "voucher": voucher})
}

func GetVoucherById(c *gin.Context) {
	var voucher models.Voucher
	id := c.Param("id")
	if err := database.DB.Where("id = ?", id).Preload("Users").Preload("Products").First(&voucher).Error; err != nil {
		c.JSON(400, gin.H{"error": "failed to get the voucher"})
		return
	}
	c.JSON(http.StatusOK, voucher)
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
	if err := database.DB.Where("id = ?", id).Preload("Users").Preload("Products").First(&voucher).Error; err != nil {
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

func GetVoucherByUserId(c *gin.Context) {
	var vouchers []models.Voucher
	user, _ := c.Get("user")
	userID := user.(models.Users).ID
	if err := database.DB.Where("user_id = ?", userID).Preload("Users").Find(&vouchers).Error; err != nil {
		c.JSON(400, gin.H{"error": "failed to get the vouchers by user id"})
		return
	}

	if len(vouchers) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "You don't have any vouchers yet"})
		return
	}

	c.JSON(http.StatusOK, vouchers)
}
