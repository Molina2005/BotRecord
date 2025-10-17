package repository

import (
	"database/sql"
	"time"
)

func QueryUser(db *sql.DB, userID int64, name string, phone string, chanel string, date time.Time) error {
	var existe bool
	// verificacion de existencia de usuario
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM usuarios WHERE id_usuario=$1)", userID).Scan(&existe)
	if err != nil {
		return err
	}

	// Si no existe se crea
	if !existe {
		_,err := db.Exec("INSERT INTO usuarios (id_usuario,nombre,telefono,canal,fecha) VALUES($1,$2,$3,$4,$5)",
			userID, name, phone, chanel, date)
		if err != nil {
			return err
		}
	}
	return nil
}

// consulta envio recordatorios
