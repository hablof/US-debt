package cache

import (
	"io"
	"os"
	"time"
)

const (
	dateFilename = "lastDate.txt"
	dateLayout   = "2006-01-02"
)

type Cache struct {
}

func (Cache) IsDateNewer(date time.Time) (bool, error) {
	dateFile, err := os.Open(dateFilename)
	if err != nil {
		return false, err
	}
	defer dateFile.Close()

	b, err := io.ReadAll(dateFile)
	if err != nil {
		return false, err
	}

	lastDate, err := time.Parse(dateLayout, string(b))
	if err != nil {
		return false, err
	}

	return date.After(lastDate), nil
}

func (Cache) UpdateDate(date time.Time) error {
	b := []byte(date.Format(dateLayout))
	return os.WriteFile(dateFilename, b, 0666)
}
