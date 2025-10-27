package apitelegram

import (
	"database/sql"
	"fmt"
	"log"
	sendmessagetelegram "modulo/SendMessageTelegram"
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
		sendmessagetelegram.MessageUser(chatID, "Uso : /registrar NombreUsuario")
	}
	dateName := partsMessage[1] // Nombre de usuario
	password := partsMessage[2]

	datePassword, err := functionsarrangements.HashPassword(password)
	if err != nil {
		sendmessagetelegram.MessageUser(chatID, "contraseña incorrecta")
		return 0, ""
	}

	date := time.Now()

	err = repository.QueryUser(db, idUser, dateName, date, datePassword)
	if err != nil {
		sendmessagetelegram.MessageUser(chatID, "error registrando usuario en la base de datos")
	} else {
		sendmessagetelegram.MessageUser(chatID, "usuario registrado correctamente")
	}
	return idUser, dateName
}

func consultIdUser(db *sql.DB, update *tgbotapi.Update, chatID int64, idUser int64) (int64, error) {
	text := update.Message.Text
	partMessage := strings.SplitN(text, " ", 2)
	if len(partMessage) < 2 {
		sendmessagetelegram.MessageUser(chatID, "Uso: /consultar contraseña")
	}
	datePassword := partMessage[1] // Contraseña

	_, err := repository.CheckUserID(db, datePassword)
	if err != nil {
		sendmessagetelegram.MessageUser(chatID, "error al buscar id principal")
	} else {
		sendmessagetelegram.MessageUser(chatID, fmt.Sprintf("id usuario: %v", idUser))
	}
	return idUser, nil
}

func deleteUser(db *sql.DB, update *tgbotapi.Update, chatID int64) int64 {
	text := update.Message.Text
	partMessage := strings.SplitN(text, " ", 2)
	if len(partMessage) < 2 {
		sendmessagetelegram.MessageUser(chatID, "Uso: /eliminar IdUsuario")
		return 0
	}
	dateId, err := strconv.ParseInt(partMessage[1], 10, 64)
	if err != nil {
		sendmessagetelegram.MessageUser(chatID, "El ID debe ser un número válido.")
		return 0
	}

	err = repository.QueryDeleteUser(db, dateId)
	if err != nil {
		sendmessagetelegram.MessageUser(chatID, "Error al eliminar usuario de la base de datos")
	} else {
		sendmessagetelegram.MessageUser(chatID, "Usuario elimiando correctamente")
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
	title := parts[3]                       // titulo recordatorio
	dateConv, err := functionsarrangements.FormatDate(dateRecord)
	if err != nil {
		sendmessagetelegram.MessageUser(chatID, "Formato de fecha incorrecto. Usa YYYY-MM-DD HH:MM")
		return
	}
	state := "pendiente"
	repeat := "no"

	errDB := repository.QueryCreateRecord(db, idUser, title, dateConv, state, repeat, channel)
	if errDB != nil {
		sendmessagetelegram.MessageUser(chatID, "error registrando recordatorio en la base de datos")
		return
	} else {
		sendmessagetelegram.MessageUser(chatID, "recordatorio registrado correctamente")
	}
}

func BotTelegram(db *sql.DB) {
	// Datos provenientes del .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Error al cargar archivo .env")
	}
	a := os.Getenv("TELEGRAM_TOKEN")
	// creacion de instancia de bot segun token
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

		// comandos para usuario
		switch {
		case strings.HasPrefix(text, "/registrar"):
			createUser(db, &update)
		case strings.HasPrefix(text, "/recordatorio"):
			createRecord(db, &update, chatID, idUser, channel)
		case strings.HasPrefix(text, "/consultar"):
			consultIdUser(db, &update, chatID, idUser)
		case strings.HasPrefix(text, "/eliminar"):
			deleteUser(db, &update, chatID)
		default:
			sendmessagetelegram.MessageUser(chatID,
				"comando no reconocido.\nusa /registrar\n/recordatorio\n/consultar\n/eliminar")
		}
	}
}
