package main

import (
	"backend/api/auth"
	"backend/api/handlers"
	"backend/api/middleware"
	"backend/initializers"
	"backend/initializers/database"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	database.InitDB()
	clearscreen()
}

func clearscreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "8080"
	}

	route := gin.Default()

	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://srmaca.vercel.app"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	route.OPTIONS("/*any", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	route.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	authGroup := route.Group("/auth")
	{
		authGroup.POST("/signup", auth.SignUp)
		authGroup.POST("/login", auth.Login)
	}
	route.GET("/validate", middleware.RequireAuth, auth.Validate)
	route.GET("/validateadmin", middleware.RequireAuth, middleware.RequireAdmin, auth.ValidateAdmin)
	orderGroup := route.Group("/order")
	orderGroup.Use(middleware.RequireAuth)
	{
		orderGroup.POST("/create", handlers.CreateOrder)
		orderGroup.GET("/get", handlers.GetOrders)
		orderGroup.GET("/getbyuserid", handlers.GetOrdersByUserID)
		orderGroup.GET("/get/:id", handlers.GetOrder)
		orderGroup.DELETE("/delete/:id", handlers.DeteleOrder)
		orderGroup.PUT("/update/:id", middleware.RequireAdmin, handlers.UpdateOrder)
	}
	voucherGroup := route.Group("/voucher")
	voucherGroup.Use(middleware.RequireAuth)
	{
		voucherGroup.POST("/create", handlers.CreateVoucher)
		voucherGroup.DELETE("/delete/:id", handlers.DeleteVoucher)
		voucherGroup.GET("/:id", handlers.GetVoucher)
		voucherGroup.GET("/", handlers.GetVouchers)
		voucherGroup.GET("/images/:id", handlers.GetVoucherImage)
	}

	route.Run("0.0.0.0:8080")
}
