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
	input := []int{1}
	_, out, err := compute.Run(buf, input)
	if err != nil {
		fmt.Printf("compute.Run: %v\n", err)
		return
	}
	fmt.Printf("out: %v\n", out)
}
