package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type parentMap map[string]string

func main() {
	content, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(content), "\n")
	parents := parentMap{}
	parentCount := map[string]int{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		entries := strings.Split(line, ")")
		parents[entries[1]] = entries[0]
		parentCount[entries[0]]++
	}
	// }
	n := 0
	for _, parent := range parents {
		n++
		for parent != "COM" {
			parent = parents[parent]
			n++
		}
	}
	fmt.Println(n)

	fmt.Println(parents.findPath("22G", "PZ4"))
	fmt.Println(parents.findPath("PZ4", "22G"))
	fmt.Println(parents.findPath("7DK", "RXL"))
	fmt.Println(parents.findPath("YOU", "SAN"))
}

func (p parentMap) findPath(a, b string) int {
	youEntries := map[string]int{}
	for i, loc := 0, a; loc != "COM"; loc = p[loc] {
		youEntries[loc] = i
		i++
	}

	foundPath := 0
	for i, loc := 0, b; loc != "COM"; loc = p[loc] {
		if j, found := youEntries[loc]; found {
			if foundPath == 0 || foundPath > i + j {
				foundPath = i + j
			}
		}
		i++
	}
	return foundPath
}
