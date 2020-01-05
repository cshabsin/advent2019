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
	for a := 0; a < 5; a++ {
		vals := []int{a}
		for b := 0; b < 5; b++ {
			if contains(vals, b) {
				continue
			}
			vals = append(vals, b)
			for c := 0; c < 5; c++ {
				if contains(vals, c) {
					continue
				}
				vals = append(vals, c)
				for d := 0; d < 5; d++ {
					if contains(vals, d) {
						continue
					}
					vals = append(vals, d)
					for e := 0; e < 5; e++ {
						if contains(vals, e) {
							continue
						}
						vals = append(vals, e)
						fmt.Printf("%v: %d\n", vals, run(buf, vals))
						vals = vals[0:4]
					}
					vals = vals[0:3]
				}
				vals = vals[0:2]
			}
			vals = vals[0:1]
		}
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
