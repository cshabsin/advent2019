package main

import (
	"fmt"
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
	for _, vals := range permutations([]int{0, 1, 2, 3, 4}, nil) {
		fmt.Printf("%v: %d\n", vals, run(buf, vals))
	}
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
		for v := range vals {
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
