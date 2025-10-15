package connect

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	// cargue de archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Error al cargar archivo .env")
	}
	// cadena de conexion para dar uso a variables de entorno
	cadenaConexion := fmt.Sprintf(
		"user=%v password=%v dbname=%v host=%v port=%v sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)
	// Abrir conexion
	db, err := sql.Open("postgres", cadenaConexion)
	if err != nil {
		return nil, err
	}
	// verificacion de conexion
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
