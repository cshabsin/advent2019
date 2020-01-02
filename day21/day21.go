package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/cshabsin/advent2019/compute"
)

var (
	walkCmds = strings.Join([]string{
		"NOT B T",
		"OR T J",
		"NOT C T",
		"OR T J",
		"AND D J", // never jump if 4 away is death
		"NOT A T", // always jump if next is death
		"OR T J",
	}, "\n") + "\nWALK\n"
	runCmds = strings.Join([]string{
		"NOT B T",
		"OR T J",
		"NOT C T",
		"OR T J", // jump if 2 or 3 away are death, unless...

		"NOT H T",
		"AND D T",
		"NOT T T",
		"AND T J", // if 8 away is death and 4 away is ok, it's a trap.

		"AND D J", // never jump if 4 away is death
		"NOT A T", // always jump if 1 away is death
		"OR T J",
	}, "\n") + "\nRUN\n"
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
	fin := make(chan bool)
	go func() {
		for {
			val, err := aIO.Read()
			if errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				fmt.Printf("read: %v\n", err)
				break
			} else if val < 256 {
				fmt.Printf("%c", rune(val))
			} else {
				fmt.Printf("Value: %d\n", val)
			}
		}
		fin <- true
	}()

	go func() {
		defer aIO.Close()
		cmds := runCmds
		for _, c := range cmds {
			if err := aIO.Write(int64(c)); err != nil {
				fmt.Printf("write: %v\n", err)
			}
			fmt.Printf("%c", c)
		}
	}()
	intcode := compute.NewIntcode(buf, bIO)
	if _, err := intcode.Run(); err != nil {
		fmt.Printf("Run: %v", err)
		return
	}
	<-fin
}
