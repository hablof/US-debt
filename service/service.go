package service

import (
	"fmt"
	"log"
	"time"

	"github.com/hablof/US-debt/model"
)

const (
	imgFilename = "output.png"

	toSoonHours = 7
	toLateHours = 22 // actualy 23
)

type Telegram interface {
	SendDebtImg(path string, diffInfo model.DebtDifference) error
	Log(message string) error
}

type Cache interface {
	GetData() (uint64, time.Time, error)
	UpdateData(data model.Debt) error
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

	// debtCandidate *model.Debt
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

	if h, _, _ := time.Now().Clock(); h >= toLateHours || h <= toSoonHours {
		log.Println("no shitpost at night please...")
		return fmt.Errorf("no shitpost at night please")
	}

	newDebt := model.Debt{}
	if err := s.prepareDebtData(&newDebt); err != nil {
		return err
	}

	if err := s.imgGen.GenerateImage(newDebt.Amuont, imgFilename); err != nil {
		fmt.Println("img generation error occured:", err)
		s.tgbot.Log(err.Error())

		return fmt.Errorf("img generation error occured: %w", err)
	}

	if err := s.tgbot.SendDebtImg(imgFilename, newDebt.Diff); err != nil {
		fmt.Println("error sending img to channel:", err)
		s.tgbot.Log(err.Error())

		return fmt.Errorf("error sending img to channel: %w", err)
	}

	if err := s.cache.UpdateData(newDebt); err != nil {
		fmt.Println("updating date error occured:", err)
		s.tgbot.Log(err.Error())
		s.cache.Erase()

		return fmt.Errorf("updating date error occured: %w", err)
	}

	return nil
}

// makes request to API US's treasury and fill unit.
// Returns error if data fetch failed.
// Also returns error if fetched date equals cached.
func (s *Service) prepareDebtData(unit *model.Debt) error {
	if err := s.ds.FetchData(); err != nil {
		log.Println("Не удалось получить данные", err)
		s.tgbot.Log("Не удалось получить данные: " + err.Error())

		return err
	}

	newDebtValue, err := s.ds.GetDebt()
	if err != nil {
		log.Println("Значение долга невалидно:", err)
		s.tgbot.Log("Значение долга невалидно: " + err.Error())

		return err
	}

	newDateValue, err := s.ds.GetDate()
	if err != nil {
		log.Println("Значение даты невалидно:", err)
		s.tgbot.Log("Значение даты невалидно:" + err.Error())

		return err
	}

	unit.Amuont = newDebtValue
	unit.Date = newDateValue

	cachedDebt, cachedDate, cacheErr := s.cache.GetData()
	switch {
	case cacheErr != nil:
		log.Println("не удалось прочитать кеш:", cacheErr)
		s.tgbot.Log("не удалось прочитать кеш: " + cacheErr.Error())
		unit.Diff.Valid = false
		s.cache.Erase()

	case unit.Date.After(cachedDate):
		log.Println("дата обновилась!")
		unit.Diff.Valid = true

		if unit.Amuont > cachedDebt {
			unit.Diff.Grows = true
			unit.Diff.Value = unit.Amuont - cachedDebt
		} else {
			unit.Diff.Grows = false
			unit.Diff.Value = cachedDebt - unit.Amuont
		}

	default:
		log.Println("дата не обновилась: не будем ничего делать")
		return fmt.Errorf("дата не обновилась: не будем ничего делать")
	}

	return nil
}
