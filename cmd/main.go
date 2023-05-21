package main

import (
	"fmt"

	debtseeker "github.com/hablof/US-debt/debt_seeker"
	imggenerator "github.com/hablof/US-debt/img_generator"
)

// const (
// 	debt = 31462676535393
// )

func main() {

	seeker := debtseeker.NewSeeker()
	debt, err := seeker.GetData()
	if err != nil {
		fmt.Println("жаль")
		return
	}

	fmt.Println(debt.GetDate())
	u, err := debt.GetDebt()
	if err != nil {
		fmt.Println("жаль 2")
		return
	}

	if err := imggenerator.GenerateImage(u); err != nil {
		fmt.Println("error occured:", err)
	}

}
