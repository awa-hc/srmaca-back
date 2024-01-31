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
	// baseURL := "http://localhost:4321/auth/verify"
	baseURL := "https://srmaca.vercel.app/auth/verify"
	verificationURL := baseURL + "?token=" + verificationToken

	// Construir el cuerpo del correo con un enlace de verificación
	subject := "Verificación de Correo Electrónico"
	htmlBody := `
		<html>
		<body style="font-family: Arial, Helvetica, sans-serif; font-size: 16px;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 4px;">
			<h1 style="font-size: 24px; color: #444;">Hola ` + username + `!</h1>
			<p style="line-height: 1.6;">
    			Gracias por registrarte en Sr Maca. Para verificar tu correo electrónico, haz clic en el siguiente botón:
			</p>
    		<div style="text-align: center;">
    			<form action="` + verificationURL + `" method="get">
        			<button style="background: #03383e; color: #fff; border: 0; padding: 12px 24px; font-size: 16px; border-radius: 4px; cursor: pointer;">Confirmar Correo</button>
    			</form>
    		</div>
    		<p style="opacity: 0.8;">
    			Este enlace expirará en 24 horas.
    		</p>
    		<p style="margin-bottom: 0;">
    			Gracias,<br>
    			El Equipo de Sr Maca
			</p>
			</div>
		</body>
		</html>
`

	// Configurar mail sender
	m := gomail.NewMessage()
	m.SetHeader("From", email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetHeader("Content-Type", "text/html; charset=UTF-8")
	m.SetBody("text/html", htmlBody)

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
