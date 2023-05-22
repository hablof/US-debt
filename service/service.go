package service

import (
	"fmt"
	"log"
	"time"
)

const (
	imgFilename = "output.png"
)

type Telegram interface {
	SendDebtImg(path string) error
	Log(message string) error
}

type Cache interface {
	IsDateNewer(date time.Time) (bool, error)
	UpdateDate(date time.Time) error
}

type DebtSeeker interface {
	FetchData() error
	GetDebt() (uint64, error)
	GetDate() (time.Time, error)
}

type ImageGenerator interface {
	GenerateImage(debt uint64, imgFilename string) error
}

type Service struct {
	tgbot  Telegram
	ds     DebtSeeker
	cache  Cache
	imgGen ImageGenerator
}

func NewService(t Telegram, ds DebtSeeker, c Cache, ig ImageGenerator) Service {
	return Service{
		tgbot:  t,
		ds:     ds,
		cache:  c,
		imgGen: ig,
	}
}

func (s *Service) Job() error {

	if err := s.ds.FetchData(); err != nil {
		log.Println("жаль", err)
		s.tgbot.Log(err.Error())
		return err
	}

	debtValue, err := s.ds.GetDebt()
	if err != nil {
		log.Println("жаль 2", err)
		s.tgbot.Log(err.Error())
		return err
	}

	debtDate, err := s.ds.GetDate()
	if err != nil {
		log.Println("жаль 3", err)
		s.tgbot.Log(err.Error())
		return err
	}

	isNewer, err := s.cache.IsDateNewer(debtDate)
	switch {
	case isNewer:
		log.Println("дата обновилась!")

	case err != nil:
		log.Println("не удалось проверить дату!", err)
		s.tgbot.Log(err.Error())

	default:
		log.Println("дата не обновилась: не будем ничего делать")
		s.tgbot.Log("дата не обновилась: не будем ничего делать")
		return fmt.Errorf("дата не обновилась: не будем ничего делать")
	}

	if err := s.imgGen.GenerateImage(debtValue, imgFilename); err != nil {
		fmt.Println("img generation error occured:", err)
		s.tgbot.Log(err.Error())
		return fmt.Errorf("img generation error occured: %w", err)
	}

	if err := s.tgbot.SendDebtImg(imgFilename); err != nil {
		fmt.Println("error sending img to channel:", err)
		s.tgbot.Log(err.Error())
		return fmt.Errorf("error sending img to channel: %w", err)
	}

	if err := s.cache.UpdateDate(debtDate); err != nil {
		fmt.Println("updating date error occured:", err)
		s.tgbot.Log(err.Error())
		return fmt.Errorf("updating date error occured: %w", err)
	}

	return nil
}
