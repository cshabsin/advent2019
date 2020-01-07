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
	var zlid, minZeroes int
	minZeroes = 1000
	for lid, l := range layers {
		zeroes := l.Count(0)
		if zeroes < minZeroes {
			minZeroes = zeroes
			zlid = lid
		}
	}
	fmt.Printf("min zeroes: %d (layer %d)\n", minZeroes, zlid)
	o := layers[zlid].Count(1)
	t := layers[zlid].Count(2)
	fmt.Printf("ones: %d, twos: %d, product: %d\n", o, t, o*t)
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

func (l layer) Count(v int) int {
	var ct int
	for _, r := range l {
		for _, c := range r {
			if c == v {
				ct++
			}
		}
	}
	return ct
}
