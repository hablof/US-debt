package main

import (
	"fmt"

	"github.com/hablof/US-debt/cache"
	debtseeker "github.com/hablof/US-debt/debt_seeker"
	imggen "github.com/hablof/US-debt/img_generator"
	"github.com/hablof/US-debt/service"
	"github.com/hablof/US-debt/tgbot"
)

func main() {

	myBot, err := tgbot.NewTgBot()
	if err != nil {
		fmt.Println("жаль, но не удалось создать бота", err)
		return
	}

	seeker := debtseeker.NewSeeker()

	s := service.NewService(myBot, seeker, cache.Cache{}, imggen.ImageGenerator{})

	fmt.Printf("s.Job(): %v\n", s.Job())
}
