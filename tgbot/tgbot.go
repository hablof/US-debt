package tgbot

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

const (
	tgChannelName  = "@usDebtDaily"
	tgChannelLogID = -1001988677061
)

type TgBot struct {
	bot *tgbotapi.BotAPI
}

func NewTgBot() (*TgBot, error) {
	godotenv.Load()
	token, found := os.LookupEnv("TOKEN")
	if !found {
		return nil, fmt.Errorf("environment variable TOKEN not found")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {

		return nil, fmt.Errorf("failed to get bot: %v", err)
	}

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	return &TgBot{bot: bot}, nil
}

func (t *TgBot) SendDebtImg(path string) error {
	myFile := tgbotapi.FilePath(path)
	pc := tgbotapi.NewPhotoToChannel(tgChannelName, myFile)
	// pc.Caption = "?"

	_, err := t.bot.Send(pc)
	return err
}

func (t *TgBot) Log(message string) error {
	mc := tgbotapi.NewMessage(tgChannelLogID, message)
	_, err := t.bot.Send(mc)
	// fmt.Printf("m.Chat.ID: %v\n", m.Chat.ID)
	return err
}
