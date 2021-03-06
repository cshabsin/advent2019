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
	io := compute.NewBufIO([]int64{1})
	intcode := compute.NewIntcode(buf, io)
	_, err = intcode.Run()
	if err != nil {
		fmt.Printf("compute.Run: %v\n", err)
		return
	}
	fmt.Printf("out(1): %v\n", io.Output())
	io = compute.NewBufIO([]int64{2})
	intcode = compute.NewIntcode(buf, io)
	_, err = intcode.Run()
	if err != nil {
		fmt.Printf("compute.Run: %v\n", err)
		return
	}
	fmt.Printf("out(2): %v\n", io.Output())
}
