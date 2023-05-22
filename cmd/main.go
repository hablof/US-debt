package main

import (
	"fmt"

	"github.com/hablof/US-debt/cache"
	debtseeker "github.com/hablof/US-debt/debt_seeker"
	imggenerator "github.com/hablof/US-debt/img_generator"
)

// const (
// 	debt = 31462676535393
// )

func main() {

	seeker := debtseeker.NewSeeker()
	dataSample, err := seeker.GetData()
	if err != nil {
		fmt.Println("жаль", err)
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

	if err := cache.UpdateDate(debtDate); err != nil {
		fmt.Println("updating date error occured:", err)
		return
	}

}
