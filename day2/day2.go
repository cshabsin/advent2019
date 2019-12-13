package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/cshabsin/advent2019/compute"
)

func main() {
	content, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	vals := strings.Split(string(content), ",")
	buf := make([]int, 0, len(vals))
	for i, v := range vals {
		intval, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			log.Fatalf("Atoi(i = %d): %v", i, err)
		}
		buf = append(buf, intval)
	}
	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			buf[1] = noun
			buf[2] = verb
			after, err := compute.Run(buf)
			if err != nil {
				fmt.Printf("error for noun %d, verb %d: %v\n", noun, verb, err)
				continue
			}
			if after[0] == 19690720 {
				fmt.Printf("noun %d, verb %d\n", noun, verb)
			}
		}
	}
}
