// email.go

package email

import (
	"crypto/tls"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gopkg.in/gomail.v2"
)

// SendVerificationEmail envía un correo de verificación al usuario con un enlace de verificación.
func SendVerificationEmail(to, username, verificationToken string) error {
	// Leer variables de entorno
	email := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")

	// Construir el cuerpo del correo con un enlace de verificación
	subject := "Verificación de Correo Electrónico"
	body := "Hola " + username + ",\n\n" +
		"Gracias por registrarte en nuestro servicio. Para verificar tu correo electrónico, haz clic en el siguiente enlace:\n\n" +
		"http://localhost:4321/auth/verify/" + verificationToken + "\n\n" +
		"Este enlace expirará en 24 horas.\n\n" +
		"Gracias,\nTu Equipo"

	// Configurar mail sender
	m := gomail.NewMessage()
	m.SetHeader("From", email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Configurar dialer
	d := gomail.NewDialer("smtp.hostinger.com", 465, email, password)
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         "smtp.hostinger.com",
	}

	// Enviar email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// GenerateVerificationToken genera un token de verificación para un usuario.
func GenerateVerificationToken(userID uint) (string, error) {
	// Define la estructura del token
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // El token expirará en 24 horas
		// Otros claims si es necesario
	}

	// Genera el token con el método de firma HMAC y la clave secreta
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// VerifyVerificationToken verifica el token de verificación y devuelve el ID del usuario.
func VerifyVerificationToken(tokenString string) (uint, error) {
	// Parsea el token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifica el método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	// Maneja errores del parsing
	if err != nil {
		return 0, err
	}

	// Verifica si el token es válido
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Verifica si el token no ha expirado
		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			if time.Now().After(expirationTime) {
				return 0, errors.New("token has expired")
			}
		}

		// Obtiene el ID del usuario desde los claims
		if sub, ok := claims["sub"].(float64); ok {
			userID := uint(sub)
			return userID, nil
		}
	}

	return 0, errors.New("invalid token")
}
