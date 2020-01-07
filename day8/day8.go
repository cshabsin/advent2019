package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	content, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	s := strings.TrimSpace(string(content))
	layers := []layer{newLayer()}
	var layerI, contentI int
	var minProd int
	minZ := 1000
	for {
		z := 0
		o := 0
		t := 0
		for y := 0; y < 6; y++ {
			for x := 0; x < 25; x++ {
				layers[layerI].Set(y, x, s[contentI])
				switch s[contentI] {
				case '0':
					z++
				case '1':
					o++
				case '2':
					t++
				}
				contentI++
			}
		}
		if z < minZ {
			minZ = z
			minProd = o*t
		}
		if contentI >= len(s) {
			break
		}
		layerI++
		layers = append(layers, newLayer())
	}
	fmt.Printf("minProd: %d\n", minProd)

	for y := 0; y < 6; y++ {
		for x := 0; x < 25; x++ {
			for _, l := range layers {
				if l[y][x] == 0 {
					fmt.Print(" ")
					break
				} else if l[y][x] == 1 {
					fmt.Print("X")
					break
				}
			}
		}
		fmt.Print("\n")
	}
}

type layer [][]int

func newLayer() layer {
	l := make(layer, 6)
	for y := range l {
		l[y] = make([]int, 25)
		for i := range l[y] {
			l[y][i] = -1
		}
	}
	return l
}

func (l layer) Set(y, x int, content byte)  {
	var v int
	switch content {
	case '0':
	case '1':
		v = 1
	case '2':
		v = 2
	default:
		fmt.Printf("invalid byte %d\n", content)
	}
	l[y][x] = v
}
