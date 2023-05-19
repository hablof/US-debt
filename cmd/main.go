package main

import (
	"fmt"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

func main() {
	img, err := imgio.Open("input.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	result := transform.ShearH(transform.ShearV(img, 30), 15)

	if err := imgio.Save("output.png", result, imgio.PNGEncoder()); err != nil {
		fmt.Println(err)
		return
	}
}
