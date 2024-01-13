package utils

import (
	"encoding/base64"
)

// DecodeBase64Image decodifica una cadena base64 y retorna los datos binarios de la imagen
func DecodeBase64Image(base64String string) ([]byte, error) {
	// Decodificar la cadena base64
	decoded, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}
