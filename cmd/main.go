package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hablof/US-debt/cache"
	debtseeker "github.com/hablof/US-debt/debt_seeker"
	imggen "github.com/hablof/US-debt/img_generator"
	"github.com/hablof/US-debt/service"
	"github.com/hablof/US-debt/tgbot"
)

const (
	interval = 4 * time.Hour
)

func main() {

	myBot, err := tgbot.NewTgBot()
	if err != nil {
		fmt.Println("жаль, но не удалось создать бота", err)
		return
	}

	seeker := debtseeker.NewSeeker()

	s := service.NewService(myBot, seeker, cache.Cache{}, imggen.ImageGenerator{})

	s.Job()

	intertuptCh := make(chan os.Signal, 1)
	signal.Notify(intertuptCh, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(interval)
	log.Println("следующая попытка в", time.Now().Add(interval).Format("15:04"))

jobLoop:
	for {
		select {
		case <-ticker.C:
			s.Job()
			log.Println("следующая попытка в", time.Now().Add(interval).Format("15:04"))

		case <-intertuptCh:
			break jobLoop
		}
	}
}
