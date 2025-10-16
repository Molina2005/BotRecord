package sendmessagetelegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

// Mensajes de exito o error que recibira el usuario
func MessageUser(ChatId int64, text string) error {
	msg := tgbotapi.NewMessage(ChatId, text)
	_, err := Bot.Send(msg)
	return err
}
// Inicializacion del bot de telegram
func Init(bot *tgbotapi.BotAPI) {
	Bot = bot
}
