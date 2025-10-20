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

		// Atributos nuevo usuario y chat de telegram
		chatID := update.Message.Chat.ID
		idUser := update.Message.From.ID
		username := update.Message.From.UserName
		phone := "0000000000"
		channel := update.Message.Chat.Type
		date := time.Now()

		err := repository.QueryUser(db, idUser, username, phone, channel, date)
		if err != nil {
			sendmessagetelegram.MessageUser(chatID, "error registrando usuario en la base de datos")
		} else {
			sendmessagetelegram.MessageUser(chatID, "usuario registrado correctamente")
		}

		// Atributos nuevo recordatorio
		text := update.Message.Text
		parts := strings.SplitN(text, " ", 3)
		if len(parts) < 3 {
			sendmessagetelegram.MessageUser(chatID, "formato incorrecto.\nUsa: YYYY-MM-DD HH:MM Título")
			return
		}

		dateRecord := parts[0] + " " + parts[1] // "2025-10-21 15:30"
		title := parts[2] // titulo del recordatorio

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
}
