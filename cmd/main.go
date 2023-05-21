package main

import (
	"fmt"

	imggenerator "github.com/hablof/US-debt/img_generator"
)

const (
	debt = 31462676535393
)

func main() {

	if err := imggenerator.GenerateImage(debt); err != nil {
		fmt.Println("error occured:", err)
	}

}
