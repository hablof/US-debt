package main

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	imgmerge "github.com/hablof/US-debt/img_merge"
	// gim "github.com/ozankasikci/go-image-merge"
)

var (
	runeToImg map[rune]string = map[rune]string{
		'0': "static/0.png",
		'1': "static/1.png",
		'2': "static/2.png",
		'3': "static/3.png",
		'4': "static/4.png",
		'5': "static/5.png",
		'6': "static/6.png",
		'7': "static/7.png",
		'8': "static/8.png",
		'9': "static/9.png",
		'_': "static/_.png",
		'$': "static/$.png",
	}
)

const (
	debt = 31462676535393

	captionWidth = 730
	captionHidth = 70

	verticalAngle   = -5
	horizontalAngle = 4
)

func main() {

	caption, err := strToRGBA(intToSpStr(debt))
	if err != nil {
		fmt.Println("cannot str to rgba", err)
		return
	}

	caption = transform.Resize(caption, captionWidth, captionHidth, transform.Lanczos)

	caption = transform.ShearH(transform.ShearV(caption, verticalAngle), horizontalAngle)

	template, err := imgio.Open("static/template.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	templateRGBA := clone.AsShallowRGBA(template)

	draw.Draw(templateRGBA, templateRGBA.Bounds(), caption, image.Point{-100, -265}, draw.Over)

	if err := imgio.Save("output.png", templateRGBA, imgio.PNGEncoder()); err != nil {
		fmt.Println(err)
		return
	}

}

func strToRGBA(str string) (*image.RGBA, error) {
	g := make([]*imgmerge.Grid, 0, len(str))

	for _, r := range str {
		g = append(g, &imgmerge.Grid{
			ImageFilePath: runeToImg[r],
		})
	}

	return imgmerge.New(g).Merge()
}

func intToSpStr(debt uint64) string {
	l := lenLoop(debt)
	outSlice := make([]byte, l+(l-1)/3)

	for i := range outSlice {
		if (i+1)%4 == 0 {
			outSlice[i] = '_'
			continue
		}

		debt, outSlice[i] = debt/10, byte(debt%10+48)
	}

	l = len(outSlice)
	for i := range outSlice {
		if i >= l-i-1 {
			break
		}
		outSlice[i], outSlice[l-i-1] = outSlice[l-i-1], outSlice[i]
	}

	return "$__" + string(outSlice)
}

func lenLoop(i uint64) int {
	if i == 0 {
		return 1
	}
	count := 0
	for i != 0 {
		i /= 10
		count++
	}
	return count
}
