package tgbot

import (
	"fmt"
	"os"
	"strings"

	"github.com/hablof/US-debt/model"

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

func (t *TgBot) SendDebtImg(path string, diffInfo model.DebtDifference) error {
	myFile := tgbotapi.FilePath(path)
	pc := tgbotapi.NewPhotoToChannel(tgChannelName, myFile)
	// pc := tgbotapi.NewPhoto(tgChannelLogID, myFile)

	if diffInfo.Valid {
		pc.Caption = formatCaption(diffInfo.Value, diffInfo.Grows)
	}

	_, err := t.bot.Send(pc)
	return err
}

func (t *TgBot) Log(message string) error {
	mc := tgbotapi.NewMessage(tgChannelLogID, message)
	_, err := t.bot.Send(mc)
	// fmt.Printf("m.Chat.ID: %v\n", m.Chat.ID)
	return err
}

// Converts uint64 to string with &nbsp as digit group separators.
// Also adds "$&nbsp" as Prefix to result.
func formatCaption(debt uint64, grow bool) string {
	l := lenLoop(debt)
	outSlice := make([]byte, l+(l-1)/3)

	for i := range outSlice {
		if (i+1)%4 == 0 {
			outSlice[i] = ' '
			continue
		}

		debt, outSlice[i] = debt/10, byte(debt%10+48)
	}

	l = len(outSlice)
	for i := range outSlice {
		if i >= l-i-1 {
			break
		}
		outSlice[i], outSlice[l-i-1] = outSlice[l-i-1], outSlice[i]
	}

	if grow {
		return "📈\u00a0" + strings.ReplaceAll("+ "+string(outSlice)+" $", " ", "\u00a0")
	}

	return "📉\u00a0" + strings.ReplaceAll("- "+string(outSlice)+" $", " ", "\u00a0")
}

func lenLoop(i uint64) int {
	if i == 0 {
		return 1
	}
	count := 0
	for i != 0 {
		i /= 10
		count++
	}
	return count
}
