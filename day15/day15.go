package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/cshabsin/advent2019/compute"
	"github.com/cshabsin/advent2019/maze"
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
	reader := bufio.NewReader(os.Stdin)
	b := &manualBoard{
		droid:  maze.MakeDroid(),
		reader: reader,
	}
	intcode := compute.NewIntcode(buf, b)
	if _, err := intcode.Run(); err != nil {
		fmt.Printf("compute.Run: %v\n", err)
		return
	}
}

type manualBoard struct {
	droid *maze.Droid

	// 1 - north
	// 2 - south
	// 3 - west
	// 4 - east
	lastMove int64
	
	reader *bufio.Reader
}

func (b *manualBoard) Read() (int64, error) {
	var text string
	for {
		fmt.Printf("input: ")
		var err error
		text, err = b.reader.ReadString('\n')
		if err != nil {
			return 0, err
		}
		text = strings.TrimSpace(text)
		if text == "n" {
			b.lastMove = 1
			return 1, nil
		} else if text == "s" {
			b.lastMove = 2
			return 2, nil
		} else if text == "w" {
			b.lastMove = 3
			return 3, nil
		} else if text == "e" {
			b.lastMove = 4
			return 4, nil
		} else if text == "p" {
			b.droid.Print()
		} else {
			i, err := strconv.ParseInt(strings.TrimSpace(text), 10, 64)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				continue
			}
			b.lastMove = i
			return i, nil
		}
	}
}

func (b *manualBoard) Write(val int64) error {
	switch b.lastMove {
	case 1: // north
		b.droid.Move(0, -1, val, true)
	case 2: // south
		b.droid.Move(0, 1, val, true)
	case 3: // west
		b.droid.Move(-1, 0, val, true)
	case 4: // east
		b.droid.Move(1, 0, val, true)
	}
	
	return nil
}
