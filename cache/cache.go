package cache

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	cacheFilename = "lastDataSample.txt"
	dateLayout    = "2006-01-02"
)

type Cache struct {
}

func (Cache) GetData() (uint64, time.Time, error) {
	dateFile, err := os.Open(cacheFilename)
	if err != nil {
		return 0, time.Time{}, err
	}
	defer dateFile.Close()

	b, err := io.ReadAll(dateFile)
	if err != nil {
		return 0, time.Time{}, err
	}

	data := strings.Fields(string(b))
	if len(data) != 2 {
		return 0, time.Time{}, fmt.Errorf("wrong data format")
	}

	lastDebt, err := strconv.Atoi(data[0])
	if err != nil {
		return 0, time.Time{}, err
	}

	lastDate, err := time.Parse(dateLayout, data[1])
	if err != nil {
		return 0, time.Time{}, err
	}

	return uint64(lastDebt), lastDate, nil
}

func (Cache) UpdateData(debt uint64, date time.Time) error {
	b := []byte(strconv.FormatUint(debt, 10) + " " + date.Format(dateLayout)) // "31462147316778 2023-05-19"
	return os.WriteFile(cacheFilename, b, 0666)
}

func (Cache) Erase() error {
	return os.Remove(cacheFilename)
}
