package main

import (
	"fmt"
)

const (
	debt = 1_023_456_789
)

func main() {
	// img, err := imgio.Open("input.png")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// result := transform.ShearH(transform.ShearV(img, -5), 8)

	// if err := imgio.Save("output.png", result, imgio.PNGEncoder()); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	fmt.Println(intToSpStr(debt))

}

func intToSpStr(debt uint64) string {
	l := lenLoop(debt)
	outSlice := make([]byte, l+(l-1)/3)

	for i := range outSlice {
		if (i+1)%4 == 0 {
			outSlice[i] = ' '
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

	return string(outSlice)
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
