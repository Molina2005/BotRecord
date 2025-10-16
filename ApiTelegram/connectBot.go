package apitelegram

import (
	"database/sql"
	// "fmt"
	"log"
	sendmessagetelegram "modulo/SendMessageTelegram"
	"modulo/repository"
	"os"
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
	a:= os.Getenv("TELEGRAM_TOKEN")
	
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
		username := update.Message.From.UserName
		telefono := "0000000000"
		canal := update.Message.Chat.Type
		fecha := time.Now()

		err := repository.QueryUser(db, chatID, username, telefono, canal, fecha)
		if err != nil {
			sendmessagetelegram.MessageUser(chatID, "error registrando usuario en la base de datos")
		}
		sendmessagetelegram.MessageUser(chatID, "usuario registrado correctamente")
	}
}
