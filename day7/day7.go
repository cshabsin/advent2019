package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/cshabsin/advent2019/compute"
)

func main() {
	content, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	buf, err := compute.ParseFile(content)
	if err != nil {
		log.Fatal(err)
	}
	var maxVals []int
	var max int64
	for _, vals := range permutations([]int{0, 1, 2, 3, 4}, nil) {
		v := run(buf, vals)
		if v > max {
			maxVals = vals
			max = v
		}
	}
	fmt.Printf("%v: %d\n", maxVals, max)

	for _, vals := range permutations([]int{5, 6, 7, 8, 9}, nil) {
		v := runFeedback(buf, vals)
		if v > max {
			maxVals = vals
			max = v
		}
	}
	fmt.Printf("%v: %d\n", maxVals, max)
}

func run(buf []int64, vals []int) int64 {
	var in int64
	for _, val := range vals {
		io := compute.NewBufIO([]int64{int64(val), in})
		intcode := compute.NewIntcode(buf, io)
		if _, err := intcode.Run(); err != nil {
			fmt.Printf("Run: %v\n", err)
			return 0
		}
		in = io.Output()[0]
	}
	return in
}

func runFeedback(buf []int64, vals []int) int64 {
	var aIOs []*compute.ChanIO

	for i, val := range vals {
		aIO, bIO := compute.NewChanIO()
		computer := compute.NewIntcode(buf, bIO)
		err := aIO.Write(int64(val))
		if err != nil {
			fmt.Printf("Write(%d): %v\n", i, err)
			return 0
		}
		aIOs = append(aIOs, aIO)
		go func(i int) {
			_, err := computer.Run()
			if err != nil {
				fmt.Printf("computers[%d].Run: %v\n", i, err)
			}
		}(i)
	}
	var in int64
	iterations := 0
	for {
		for i := range vals {
			err := aIOs[i].Write(in)
			if err != nil {
				fmt.Printf("aIOs[%d].Write(%d): %v\n", i, err)
				return 0
			}
			newIn, err := aIOs[i].Read()
			if errors.Is(err, io.EOF) {
				return in
			}
			if err != nil {
				fmt.Printf("aIOs[%d].Read(%d): %v\n", i, err)
				return 0
			}
			in = newIn
		}
		iterations++
	}

}

func contains(vals []int, v int) bool {
	for _, i := range vals {
		if i == v {
			return true
		}
	}
	return false
}

func makeArr(arr []int, v int) []int {
	out := make([]int, 0, len(arr)+1)
	for _, a := range arr {
		out = append(out, a)
	}
	out = append(out, v)
	return out
}

func permutations(vals []int, in [][]int) [][]int {
	if in == nil {
		in = [][]int{}
		for _, v := range vals {
			in = append(in, []int{v})
		}
	}
	out := [][]int{}
	for _, inVal := range in {
		for _, v := range vals {
			if contains(inVal, v) {
				continue
			}
			out = append(out, makeArr(inVal, v))
		}
	}
	if len(out[0]) == len(vals) {
		return out
	}
	return permutations(vals, out)
}
