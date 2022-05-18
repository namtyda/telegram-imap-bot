package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/mail"
)

const (
	commandStart      = "start"
	commandAddMail    = "add_mail"
	commandDeleteMail = "delete_mail"
	commandSecret     = "secret"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	case commandAddMail:
		return b.handleAddMailCommand(message)
	case commandDeleteMail:
		return b.handleDeleteMailCommand(message)
	case commandSecret:
		return b.handleSecretCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Command())
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Undefined command")
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "it's only message "+message.Text)
	_, err := b.bot.Send(msg)
	return err
}
func (b *Bot) handleAddMailCommand(message *tgbotapi.Message) error {
	addr, err := mail.ParseAddress(message.CommandArguments())

	if err != nil {
		return invalidEmail
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, addr.Address)

	_, errSend := b.bot.Send(msg)
	return errSend
}

func (b *Bot) handleDeleteMailCommand(message *tgbotapi.Message) error {
	addr, err := mail.ParseAddress(message.CommandArguments())

	if err != nil {
		return invalidEmail
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, addr.Address)

	_, errSend := b.bot.Send(msg)
	return errSend
}

func (b *Bot) handleSecretCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Chat.FirstName+" Жопа")
	_, err := b.bot.Send(msg)
	return err
}
