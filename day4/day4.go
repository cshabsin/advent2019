package main

import (
	"fmt"
	"strconv"
)

func match(in int) bool {
	s := strconv.Itoa(in)
	prev := 0
	dupCount := 0
	var foundDup bool
	for i := 0; i < len(s); i++ {
		c := int(s[i]-'0')
		if c == prev {
			dupCount++
		} else if c < prev {
			return false
		} else {
			if dupCount == 1 {
				foundDup = true
			}
			dupCount = 0
		}
		prev = c
	}
	if dupCount == 1 {  // TODO: yuck
		foundDup = true
	}
	return foundDup
}

func main() {
	var c int
	for i := 382345; i <= 843167; i++ {
		if match(i) {
			c++
		}
	}
	fmt.Printf("count: %d\n", c)
}
