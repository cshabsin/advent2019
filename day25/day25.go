package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

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
	aIO.NoTimeout = true
	bIO.NoTimeout = true
	
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

	intcode := compute.NewIntcode(buf, bIO)
	go func() {
		if _, err := intcode.Run(); err != nil {
			fmt.Printf("Run: %v", err)
		}
	}()
	defer aIO.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		move, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("ReadString: %v", err)
			break
		}
		for _, r := range move {
			if err := aIO.Write(int64(r)); err != nil {
				fmt.Printf("Write: %v", err)
				break
			}
		}
	}
}
