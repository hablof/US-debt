package main

import (
	"fmt"

	"github.com/hablof/US-debt/cache"
	debtseeker "github.com/hablof/US-debt/debt_seeker"
	imggenerator "github.com/hablof/US-debt/img_generator"
	"github.com/hablof/US-debt/tgbot"
)

const (
	imgFilename = "output.png"
)

func main() {

	myBot, err := tgbot.NewTgBot()
	if err != nil {
		fmt.Println("жаль, но не удалось создать бота", err)
		return
	}

	seeker := debtseeker.NewSeeker()
	dataSample, err := seeker.GetData()
	if err != nil {
		fmt.Println("жаль", err)
		fmt.Printf("myBot.ReportErrInfo(err): %v\n", myBot.Log(err.Error()))
		return
	}

	debtValue, err := dataSample.GetDebt()
	if err != nil {
		fmt.Println("жаль 2", err)
		return
	}

	debtDate, err := dataSample.GetDate()
	if err != nil {
		fmt.Println("жаль 3", err)
		return
	}

	isNewer, err := cache.IsDateNewer(debtDate)
	switch {
	case isNewer:
		fmt.Println("дата обновилась!")

	case err != nil:
		fmt.Println("не удалось проверить дату!", err)

	default:
		fmt.Println("дата не обновилась: не будем ничего делать")
		return
	}

	if err := imggenerator.GenerateImage(debtValue); err != nil {
		fmt.Println("img generation error occured:", err)
		return
	}

	if err := myBot.SendDebtImg(imgFilename); err != nil {
		fmt.Println("error sending img to channel:", err)
		return
	}

	if err := cache.UpdateDate(debtDate); err != nil {
		fmt.Println("updating date error occured:", err)
		return
	}

}
