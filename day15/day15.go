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
	b := &board{
		reader: reader,
	}
	b.setVal(0, 0, 1)  // empty
	intcode := compute.NewIntcode(buf, b)
	if _, err := intcode.Run(); err != nil {
		fmt.Printf("compute.Run: %v\n", err)
		return
	}
}

type board struct {
	posX, posY int

	// 1 - north
	// 2 - south
	// 3 - west
	// 4 - east
	lastMove int64
	
	board maze.Board

	reader *bufio.Reader
}

func (b *board) Read() (int64, error) {
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
			b.board.Print()
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

func (b *board) setVal(posX, posY, val int) {
	b.board.SetVal(posX, posY, val)
}

func (b *board) Write(val int64) error {
	nextPosX := b.posX
	nextPosY := b.posY
	switch b.lastMove {
	case 1: // north
		nextPosY--
	case 2: // south
		nextPosY++
	case 3: // west
		nextPosX--
	case 4: // east
		nextPosX++
	}
	switch val {
	case 0:
		fmt.Printf("wall at %d, %d; no move\n", nextPosX, nextPosY)
		b.setVal(nextPosX, nextPosY, 2)
	case 1:
		fmt.Printf("moved to %d, %d\n", nextPosX, nextPosY)
		b.setVal(nextPosX, nextPosY, 1)
		b.posX = nextPosX
		b.posY = nextPosY
	case 2:
		fmt.Printf("oxygen at %d, %d; moved in\n", nextPosX, nextPosY)
		b.setVal(nextPosX, nextPosY, 3)
		b.posX = nextPosX
		b.posY = nextPosY
	}
	
	return nil
}
