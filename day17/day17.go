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

	aIO, bIO := compute.NewChanIO()
	go func() {
		for {
			val, err := aIO.Read()
			if err != nil {
				fmt.Printf("aIO.Read: %v", err)
				return
			}
			fmt.Printf("%c", rune(val))
		}
	}()
	intcode := compute.NewIntcode(buf, bIO)
	if _, err := intcode.Run(); err != nil {
		fmt.Printf("compute.Run: %v\n", err)
		return
	}
	
}
