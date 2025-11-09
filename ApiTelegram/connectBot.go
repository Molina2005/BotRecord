package apitelegram

import (
	"database/sql"
	"fmt"
	"log"
	sendmessagetelegram "modulo/SendMessageTelegram"
	"modulo/bot"
	functionsarrangements "modulo/functionsArrangements"
	"modulo/repository"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func createUser(db *sql.DB, update *tgbotapi.Update) (int64, string) {
	// Datos automaticos y digitados por usuario
	chatID := update.Message.Chat.ID
	idUser := update.Message.From.ID
	text := update.Message.Text
	partsMessage := strings.SplitN(text, " ", 3)

	if len(partsMessage) < 3 {
		sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("%v Uso : /registrar NombreUsuario", bot.EmojiAtencion))
		return 0, ""
	}
	dateName := partsMessage[1] // Nombre usuario
	password := partsMessage[2] // Contraseña

	// Hasheo(encriptacion) contraseña usuario
	datePassword, err := functionsarrangements.HashPassword(password)
	if err != nil {
		sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("%v contraseña incorrecta", bot.EmojiAdvertencia))
		return 0, ""
	}
	date := time.Now()

	err = repository.QueryUser(db, idUser, dateName, date, datePassword)
	if err != nil {
		sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("%v error registrando usuario en la base de datos", bot.EmojiError))
		return 0, ""
	} else {
		sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("%v usuario registrado correctamente", bot.EmojiExito))
	}
	return idUser, dateName
}

func consultIdUser(db *sql.DB, update *tgbotapi.Update, chatID int64) (int64, error) {
	text := update.Message.Text
	partMessage := strings.SplitN(text, " ", 2)
	if len(partMessage) < 2 {
		sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("%v Uso: /consultar contraseña", bot.EmojiAtencion))
		return 0, nil
	}
	datePassword := partMessage[1] // Contraseña

	idUser, err := repository.CheckUserID(db, datePassword, chatID)
	if err != nil {
		sendmessagetelegram.MessageUser(chatID, err.Error())
		return 0, err
	} else {
		sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("%v Codigo unico: %v", bot.EmojiId, idUser))
	}
	return idUser, nil
}

func deleteUser(db *sql.DB, update *tgbotapi.Update, chatID int64) int64 {
	text := update.Message.Text
	partMessage := strings.SplitN(text, " ", 2)
	if len(partMessage) < 2 {
		sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("%v Uso: /eliminar IdUsuario", bot.EmojiAtencion))
		return 0
	}
	dateId, err := strconv.ParseInt(partMessage[1], 10, 64)
	if err != nil {
		sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("%v El ID debe ser un número válido.", bot.EmojiAdvertencia))
		return 0
	}

	err = repository.QueryDeleteUser(db, dateId)
	if err != nil {
		sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("%v Error al eliminar usuario de la base de datos", bot.EmojiError))
		return 0
	} else {
		sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("%v Usuario elimiando correctamente", bot.EmojiEliminacion))
	}
	return dateId
}

func createRecord(db *sql.DB, update *tgbotapi.Update, chatID, idUser int64, channel string) {
	text := update.Message.Text
	parts := strings.SplitN(text, " ", 4)
	if len(parts) < 4 {
		sendmessagetelegram.MessageUser(chatID, "formato incorrecto.\nUsa: YYYY-MM-DD HH:MM Título")
		return
	}
	dateRecord := parts[1] + " " + parts[2] // "2025-10-21 15:30"
	title := parts[3]
	dateConv, err := functionsarrangements.FormatDate(dateRecord)
	if err != nil {
		sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("%v Formato de fecha incorrecto. Usa YYYY-MM-DD HH:MM", bot.EmojiAtencion))
		return
	}
	state := "pendiente"

	errDB := repository.QueryCreateRecord(db, idUser, title, dateConv, state, channel)
	if errDB != nil {
		sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("%v error registrando recordatorio en la base de datos", bot.EmojiError))
		return
	} else {
		sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("%v recordatorio registrado correctamente", bot.EmojiExito))
	}
}

func sendReminder(db *sql.DB, chatID int64) {
	for {
		reminders, err := repository.ConsultShippingReminder(db, chatID)
		if err != nil {
			fmt.Println("Error al consultar recordatorio", err)
			return
		}
		loc, _ := time.LoadLocation("America/Bogota")
		currentDate := time.Now().In(loc).Truncate(time.Minute) // sin segundos

		for _, r := range reminders {
			recordDate := r.DateRecord.Truncate(time.Minute)
			if recordDate.Equal(currentDate) {
				formatDate := r.DateRecord.Format("15:04")
				sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("%v %v %v", bot.EmojiNotificacion, r.Title, formatDate))
				_, err := db.Exec("UPDATE recordatorios SET estado = 'enviado' WHERE id_usuario = $1 AND id_recordatorio = $2", int64(chatID), int64(r.IdRecordatorios))
				if err != nil {
					log.Println("Error actualizando el estado del recordatorio", err)
				}
			}
		}
		// lapso de 1min para que vuelva repetir el proceso
		time.Sleep(1 * time.Minute)
	}
}

func consultReminder(db *sql.DB, update *tgbotapi.Update, ChatID int64) {
	text := update.Message.Text
	parts := strings.Split(text, " ")
	if len(parts) < 1 {
		sendmessagetelegram.MessageUser(ChatID, fmt.Sprintf("%v Uso: /Lista", bot.EmojiAtencion))
	}
	reminders, err := repository.ConsultShippingReminder(db, ChatID)
	if err != nil {
		return
	}

	for _, r := range reminders {
		recordDate := r.DateRecord.Format("2006-01-02 15:04")
		sendmessagetelegram.MessageUser(ChatID, fmt.Sprintf("%v [%v] %v %v %v ", bot.EmojiRecordatorio, r.IdRecordatorios, r.Title, recordDate, r.Estado))
	}
}

func deleteReminders(db *sql.DB, update *tgbotapi.Update, chatID int64) {
	text := update.Message.Text
	parts := strings.SplitN(text, " ", 2)
	if len(parts) < 2 {
		sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("%v Uso: /eliminarrecordatorio codigorecordatorio", bot.EmojiAtencion))
	}
	partId := parts[1]
	Id, _ := strconv.Atoi(partId)

	rows, err := repository.DeleteReminder(db, Id, chatID)
	if err != nil {
		sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("%v Error al eliminar recordatorio", bot.EmojiError))
	} else if rows == 0 {
		sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("%v No se encontró ningún recordatorio con ese ID", bot.EmojiAdvertencia))
	} else {
		sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("%v Recordatorio eliminado correctamente", bot.EmojiExito))
	}
}

func BotTelegram(db *sql.DB) {
	// Datos provenientes del .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Error al cargar archivo .env")
	}
	a := os.Getenv("TELEGRAM_TOKEN")
	// Creacion de instancia de bot segun token
	bot, err := tgbotapi.NewBotAPI(a)
	if err != nil {
		log.Println("Error al crear el bot:", err)
	}
	// Nombre de usuario del bot
	log.Printf("Bot autorizado como: %v", bot.Self.UserName)
	// Guarda la instancia del bot para siempre trabajar con el mismo bot
	sendmessagetelegram.Init(bot)
	// Estructura para poder recibir mensajes nuevos
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	// Creacion canal de actualizaciones
	updates := bot.GetUpdatesChan(u)

	// Se procesa cada mensaje recibido
	for update := range updates {
		if update.Message == nil {
			continue
		}
		// Datos para envio de recordatorio
		text := update.Message.Text
		idUser := update.Message.From.ID
		chatID := update.Message.Chat.ID
		channel := update.Message.Chat.Type

		// Opciones de comandos para usuario
		msg := `Usa alguno de los siguientes comandos:

		❗Recomendacion: hacer uso de espacios entre apartados
		
		📋 /menu 
		→ Desplegar menu de comandos

		📝 /registrar nombreusuario contraseña
		→ Registra un nuevo usuario.

		⏰ /recordatorio fechayhora descripción
		Ejemplo: /recordatorio 2025-10-25 14:30 Reunión 

		❌ /borrarrecordatorio id
		→ Borrar recordatorio de la lista

		🔍 /consultar contraseña
		→ Consulta tu ID de usuario.

		🗑️ /eliminar id	usuario
		→ Elimina tu cuenta del sistema.

		🗓️ /lista 
		→ Lista de recordatorios`

		// Comandos para usuario
		switch {
		case strings.HasPrefix(text, "/menu"):
			sendmessagetelegram.MessageUser(chatID, msg)
		case strings.HasPrefix(text, "/registrar"):
			createUser(db, &update)
		case strings.HasPrefix(text, "/recordatorio"):
			createRecord(db, &update, chatID, idUser, channel)
		case strings.HasPrefix(text, "/consultar"):
			consultIdUser(db, &update, chatID)
		case strings.HasPrefix(text, "/eliminar"):
			deleteUser(db, &update, chatID)
		case strings.HasPrefix(text, "/lista"):
			consultReminder(db, &update, chatID)
		case strings.HasPrefix(text, "/borrarrecordatorio"):
			deleteReminders(db, &update, chatID)
		default:
			sendmessagetelegram.MessageUser(chatID, msg)
		}
		// Envio recordatorios al usurio en fecha especifica
		go sendReminder(db, chatID)
	}
}
