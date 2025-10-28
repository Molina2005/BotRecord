package repository

import (
	"database/sql"
	sendmessagetelegram "modulo/SendMessageTelegram"
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
func CheckUserID(db *sql.DB, Inputpassword string, ChatID int64) (int64, error) {
	var hashedPassword string
	// Se valida que sea el id del usuario que digito el comando
	err := db.QueryRow("SELECT id_usuario, contrasena FROM usuarios WHERE id_usuario = $1", ChatID).Scan(&ChatID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			sendmessagetelegram.MessageUser(ChatID, "Usuario no encontrado. Regístrate primero con /registrar.")
			return 0, nil
		}
		return 0, nil
	}

	// Comparacion contraseña en base datos y digitada por usuario en Telegram
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(Inputpassword))
	if err != nil {
		sendmessagetelegram.MessageUser(ChatID, "Contraseña incorrecta")
		return 0, nil
	}
	return ChatID, nil
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

// Consulta envio recordatorio
// func ConsultShippingReminder(db *sql.DB, ChatID int64, title string, date_record time.Time) (int64, error) {

// 	rows, err := db.Query("SELECT id_usuario, titulo, fecha_recordatorio FROM recordatorios WHERE id_usuario = $1", ChatID)
// 	if err != nil {
// 		return 0, nil
// 	}
// 	return ChatID, nil
// }
