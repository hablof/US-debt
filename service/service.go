package service

import (
	"fmt"
	"log"
	"time"

	"github.com/hablof/US-debt/model"
)

const (
	imgFilename = "output.png"
)

type Telegram interface {
	SendDebtImg(path string, diffInfo model.DebtDifference) error
	Log(message string) error
}

type Cache interface {
	GetData() (uint64, time.Time, error)
	UpdateData(debt uint64, date time.Time) error
	Erase() error
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

	newDebt, err := s.ds.GetDebt()
	if err != nil {
		log.Println("жаль 2", err)
		s.tgbot.Log(err.Error())

		return err
	}

	newDate, err := s.ds.GetDate()
	if err != nil {
		log.Println("жаль 3", err)
		s.tgbot.Log(err.Error())

		return err
	}

	diffInfo := model.DebtDifference{}

	cachedDebt, cachedDate, cacheErr := s.cache.GetData()
	switch {
	case cacheErr != nil:
		log.Println("не удалось проверить дату!", cacheErr)
		s.tgbot.Log(cacheErr.Error())
		diffInfo.Valid = false
		s.cache.Erase()

	case newDate.After(cachedDate):
		log.Println("дата обновилась!")
		diffInfo.Valid = true
		if newDebt > cachedDebt {
			diffInfo.Grows = true
			diffInfo.Value = newDebt - cachedDebt
		} else {
			diffInfo.Grows = false
			diffInfo.Value = cachedDebt - newDebt
		}

	default:
		log.Println("дата не обновилась: не будем ничего делать")
		s.tgbot.Log("дата не обновилась: не будем ничего делать")

		return fmt.Errorf("дата не обновилась: не будем ничего делать")
	}

	if err := s.imgGen.GenerateImage(newDebt, imgFilename); err != nil {
		fmt.Println("img generation error occured:", err)
		s.tgbot.Log(err.Error())

		return fmt.Errorf("img generation error occured: %w", err)
	}

	if err := s.tgbot.SendDebtImg(imgFilename, diffInfo); err != nil {
		fmt.Println("error sending img to channel:", err)
		s.tgbot.Log(err.Error())

		return fmt.Errorf("error sending img to channel: %w", err)
	}

	if err := s.cache.UpdateData(newDebt, newDate); err != nil {
		fmt.Println("updating date error occured:", err)
		s.tgbot.Log(err.Error())
		s.cache.Erase()

		return fmt.Errorf("updating date error occured: %w", err)
	}

	return nil
}
