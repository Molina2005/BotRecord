package repository

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func QueryUser(db *sql.DB, userID int64, name string, date time.Time, password string) error {
	var existe bool
	// Verificacion de existencia de usuario
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM usuarios WHERE id_usuario=$1)", userID).Scan(&existe)
	if err != nil {
		return err
	}

	// Si no existe se crea
	if !existe {
		_, err := db.Exec("INSERT INTO usuarios (id_usuario,nombre,fecha,contrasena) VALUES($1,$2,$3,$4)",
			userID, name, date, password)
		if err != nil {
			return err
		}
	}
	return nil
}

// Consultar id usuario
func CheckUserID(db *sql.DB, Inputpassword string) (int64, error) {
	var id_user int64
	var hashedPassword string
	err := db.QueryRow("SELECT id_usuario, contrasena FROM usuarios").Scan(&id_user, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("no se encontro usuario con ese id")
			return 0, nil
		}
		return 0, nil
	}

	// Comparacion contraseña en base datos y digitada por usuario en Telegram
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(Inputpassword))
	if err != nil {
		fmt.Println("Contraseña incorrecta")
		return 0, nil
	}

	return id_user, nil
}

// Consulta eliminar usuario
func QueryDeleteUser(db *sql.DB, userID int64) error {
	_, err := db.Exec("DELETE FROM usuarios WHERE id_usuario=$1", userID)
	return err
}

// Consulta insertar recordatorios
func QueryCreateRecord(
	db *sql.DB, id_user int64,
	title string, date_record time.Time,
	state, shipping_chanel string,
) error {
	_, err := db.Exec(
		"INSERT INTO recordatorios (id_usuario,titulo,fecha_recordatorio,estado,canal_envio,fecha_creacion) VALUES($1,$2,$3,$4,$5, NOW())",
		id_user, title, date_record, state, shipping_chanel)
	if err != nil {
		return err
	}
	return nil
}
