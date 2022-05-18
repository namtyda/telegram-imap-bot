package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	invalidEmail = errors.New("invalid mail")
)

func (b *Bot) handleError(chatID int64, err error) {
	var messageText string

	switch err {
	case invalidEmail:
		messageText = "mail is invalid, pass valid email into command"
	default:
		messageText = "Default error"
	}

	msg := tgbotapi.NewMessage(chatID, messageText)
	b.bot.Send(msg)
}
