package middleware

import (
	"backend/initializers/database"
	"backend/initializers/models"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

/* func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Auth")
	if err != nil {
		fmt.Println("Error obteniendo cookie 'Auth':", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	fmt.Println("Token de la cookie 'Auth':", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			fmt.Println("Token expirado")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user models.Users
		if err := database.DB.First(&user, claims["sub"]).Error; err != nil {
			fmt.Println("Clave 'sub' no encontrada en los claims")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", user)
		c.Next()
	} else {
		fmt.Println("Token no válido")
		c.AbortWithStatus(http.StatusUnauthorized)
	}
} */

func RequireAuth(c *gin.Context) {
	// Obtener el token del encabezado "Authorization"
	tokenString, err := c.Cookie("Auth")
	if err != nil {
		fmt.Println("Error obteniendo cookie 'Auth':", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			fmt.Println("Token expirado")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user models.Users
		if err := database.DB.First(&user, claims["sub"]).Error; err != nil {
			fmt.Println("Clave 'sub' no encontrada en los claims")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", user)
		c.Next()
	} else {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("Error: Header 'Authorization' no presente")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// El token debe tener el formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Println("Error: Formato de token inválido")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				fmt.Println("Token expirado")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			var user models.Users
			if err := database.DB.First(&user, claims["sub"]).Error; err != nil {
				fmt.Println("Clave 'sub' no encontrada en los claims")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			c.Set("user", user)
			c.Next()
		} else {
			fmt.Println("Token no válido")
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func RequireAdmin(c *gin.Context) {
	tokenString, err := c.Cookie("Auth")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		var user models.Users
		database.DB.First(&user, claims["sub"])

		if token.Claims.(jwt.MapClaims)["role"] != "admin" {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("user", user)
		c.Next()
	} else {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("Error: Header 'Authorization' no presente")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// El token debe tener el formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Println("Error: Formato de token inválido")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				fmt.Println("Token expirado")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			var user models.Users
			if err := database.DB.First(&user, claims["sub"]).Error; err != nil {
				fmt.Println("Clave 'sub' no encontrada en los claims")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			database.DB.First(&user, claims["sub"])

			if token.Claims.(jwt.MapClaims)["role"] != "admin" {
				c.AbortWithStatus(http.StatusUnauthorized)
			}
			c.Set("user", user)
			c.Next()
		} else {
			fmt.Println("Token no válido")
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
