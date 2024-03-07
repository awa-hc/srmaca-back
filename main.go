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
		// AllowOrigins: []string{"http://localhost:3001"},
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
		authGroup.GET("/verify/:token", handlers.VerifyEmail)
		authGroup.POST("/forgotPassword", auth.ForgotPassword)
		authGroup.POST("/resetPassword", auth.ResetPassword)
	}
	route.GET("/validate", middleware.RequireAuth, auth.Validate)
	route.GET("/validateadmin", middleware.RequireAuth, auth.ValidateAdmin)

	voucherGroup := route.Group("/voucher")
	voucherGroup.Use(middleware.RequireAuth)
	{

		voucherGroup.POST("/create", handlers.CreateVoucher)
		voucherGroup.PUT("/confirm/:id", handlers.ConfirmVoucher)
		voucherGroup.DELETE("/delete/:id", handlers.DeleteVoucher)
		voucherGroup.GET("/id/:id", handlers.GetVoucherById)
		voucherGroup.GET("/:id", handlers.GetVoucher)
		voucherGroup.GET("/", handlers.GetVouchers)
		voucherGroup.GET("/images/:id", handlers.GetVoucherImage)
		voucherGroup.GET("/user/", handlers.GetVoucherByUserId)
	}

	userGroup := route.Group("/user")
	userGroup.Use(middleware.RequireAuth)
	{
		userGroup.GET("/", handlers.GetUser)
		userGroup.GET("/:id", handlers.GetUserById)
	}
	productGroup := route.Group("/product")
	productGroup.Use(middleware.RequireAuth)
	{
		productGroup.POST("/create", handlers.CreateProduct)
		productGroup.GET("/", handlers.GetProducts)
		productGroup.GET("/:id", handlers.GetProduct)
		productGroup.PUT("/update/:id", handlers.UpdateProduct)
	}

	emailGroup := route.Group("/email")
	{
		emailGroup.POST("/contact", handlers.ContactEmail)
	}

	// route.Run(":" + Port)
	route.Run("0.0.0.0:" + Port)
}
