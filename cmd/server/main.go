package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"gitlab.ozon.dev/namtyda/homework-2/configs"
	"gitlab.ozon.dev/namtyda/homework-2/internal/app/imap"
	"gitlab.ozon.dev/namtyda/homework-2/internal/app/telegram"
)

func init() {
	if err := configs.InitConfig(); err != nil {
		log.Fatalf("Error init configs %s", err.Error())
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(viper.GetString("apiKeys.telegram"))
	bot.Debug = true

	if err != nil {
		log.Fatalf("Error init tgbotapi  %s", err.Error())
	}
	tg := telegram.New(bot)

	go (func() {
		if err := tg.Start(); err != nil {
			log.Fatalf("Error start bot %s", err.Error())
		}
	})()

	mailClient := imap.NewClient("host:993", "username", "password")
	if errLogin := mailClient.Login(); errLogin != nil {
		log.Fatalf("Login imap erros %s", errLogin.Error())
	}
	defer mailClient.Logout()

	sendMsg := func(msg *imap.Message) {
		text := fmt.Sprintf("*%s*", msg.Subject)
		tg.SendMessage(text, 23223)
		mailClient.MarkMsgSeen(msg)
	}

	for _, msg := range mailClient.FetchUnseenMsgs() {
		sendMsg(msg)
	}

	msgs := make(chan *imap.Message)

	go func() {
		for msg := range msgs {
			sendMsg(msg)
		}
	}()

	mailClient.WaitNewMsgs(msgs, time.Duration(2*time.Second))
}
