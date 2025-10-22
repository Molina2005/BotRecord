package apitelegram

import (
	"database/sql"
	"log"
	sendmessagetelegram "modulo/SendMessageTelegram"

	functionsarrangements "modulo/functionsArrangements"
	"modulo/repository"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func CreateUser(db *sql.DB, update *tgbotapi.Update) (int64, string) {
	// Atributos nuevo usuario y chat de telegram
	chatID := update.Message.Chat.ID
	idUser := update.Message.From.ID
	username := update.Message.Text
	partsMessage := strings.SplitN(username, " ", 3)
	dateName := partsMessage[0]
	conv := strings.ReplaceAll(dateName, ",", "")
	phone := "0000000000"
	channel := update.Message.Chat.Type
	date := time.Now()

	err := repository.QueryUser(db, idUser, conv, phone, channel, date)
	if err != nil {
		sendmessagetelegram.MessageUser(chatID, "error registrando usuario en la base de datos")
	} else {
		sendmessagetelegram.MessageUser(chatID, "usuario registrado correctamente")
	}
	return idUser, channel
}

func CreateRecord(db *sql.DB, update *tgbotapi.Update, idUser int64, channel string) {
	// Atributos nuevo recordatorio
	chatID := update.Message.Chat.ID
	text := update.Message.Text
	parts := strings.SplitN(text, " ", 3)
	if len(parts) < 3 {
		sendmessagetelegram.MessageUser(chatID, "formato incorrecto.\nUsa: YYYY-MM-DD HH:MM Título")
		return
	}
	dateRecord := parts[0] + " " + parts[1] // "2025-10-21 15:30"
	title := parts[2]                       // titulo del recordatorio
	dateConv, err := functionsarrangements.FormatDate(dateRecord)
	if err != nil {
		sendmessagetelegram.MessageUser(chatID, "Formato de fecha incorrecto. Usa YYYY-MM-DD HH:MM")
		return
	}
	state := "pendiente"
	repeat := "no"
	timeRecord := time.Now()

	// envio recordatorios
	errDB := repository.CreateRecord(db, idUser, title, state, repeat, channel, timeRecord, dateConv)
	if errDB != nil {
		sendmessagetelegram.MessageUser(chatID, "error registrando recordatorio en la base de datos")
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
	// Cadena de conexion con variables de entorno
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
		CreateUser(db, &update)
		idUser := update.Message.From.ID
		channel := update.Message.Chat.Type
		CreateRecord(db, &update, idUser, channel)
	}
}
