package repository

import (
	"database/sql"
	// "fmt"
	"time"
)

func QueryUser(db *sql.DB, userID int64, name string, date time.Time, password string) error {
	var existe bool
	// verificacion de existencia de usuario
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

// consultar id usuario
// func CheckUserID(db *sql.DB, userID int64) (int64, error) {
// 	var id_user int64
// 	err := db.QueryRow("SELECT u.id_usuario FROM usuarios u WHERE id_usuario = $1", userID).Scan(&id_user)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			fmt.Println("no se encontro usuario con ese id")
// 			return 0, nil
// 		}
// 		return 0, nil
// 	}
// 	return id_user, nil
// }

// consulta eliminar usuario
func QueryDeleteUser(db *sql.DB, userID int64) error {
	_, err := db.Exec("DELETE FROM usuarios WHERE id_usuario=$1", userID)
	return err
}

// consulta insertar recordatorios
func QueryCreateRecord(
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
