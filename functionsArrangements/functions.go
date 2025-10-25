package functionsarrangements

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Formateo de fechas y hora
func FormatDate(date string) (time.Time, error) {
	formatDate := "2006-01-02 15:04"
	loc, _ := time.LoadLocation("America/Bogota")
	return time.ParseInLocation(formatDate, date, loc)
}

// Cifrado de contraseñas
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
