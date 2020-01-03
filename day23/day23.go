package main

import (
	"io/ioutil"
	"log"

	"github.com/cshabsin/advent2019/compute"
	"github.com/cshabsin/advent2019/network"
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
	net := network.NewNetwork()
	for i := 0; i < 50; i++ {
		net.AddComputer(buf, i)
	}
	net.Run()
}
