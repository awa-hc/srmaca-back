package handlers

import (
	"backend/initializers/database"
	"backend/initializers/models"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func GetVoucherImage(c *gin.Context) {
	var voucher models.Voucher
	voucherID := c.Param("id")

	if err := database.DB.Where("id = ?", voucherID).First(&voucher).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "voucher not found"})
		return
	}

	voucherPath := filepath.Join("uploads", voucher.Img)
	imageData, err := os.ReadFile(voucherPath)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read image"})
		return
	}
	c.Data(http.StatusOK, "image/png", imageData)
}
