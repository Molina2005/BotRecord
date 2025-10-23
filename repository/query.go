package repository

import (
	"database/sql"
	"time"
)

func QueryUser(db *sql.DB, userID int64, name string, date time.Time) error {
	var existe bool
	// verificacion de existencia de usuario
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM usuarios WHERE id_usuario=$1)", userID).Scan(&existe)
	if err != nil {
		return err
	}

	// Si no existe se crea
	if !existe {
		_, err := db.Exec("INSERT INTO usuarios (id_usuario,nombre,fecha) VALUES($1,$2,$3)",
			userID, name, date)
		if err != nil {
			return err
		}
	}
	return nil
}

// consulta insertar recordatorios
func CreateRecord(
	db *sql.DB, id_user int64,
	title string, date_record time.Time,
	state, repeat, shipping_chanel string,
) error {
	_, err := db.Exec(
		"INSERT INTO recordatorios (id_usuario,titulo,fecha_recordatorio,estado,repetir,canal_envio,fecha_creacion) VALUES($1,$2,$3,$4,$5,$6, NOW())",
		id_user, title, date_record, state, repeat, shipping_chanel)
	if err != nil {
		return err
	}

	return nil
}
