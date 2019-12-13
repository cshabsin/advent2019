package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"./fuel"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	cumulative := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		mass,err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		cumulative += fuel.AllFuel(mass)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Cumulative Mass: %d\n", cumulative)
}
