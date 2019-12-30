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
	init := 50
	b := &board{
		posX: init,
		posY: init,
		minX: init,
		minY: init,
		maxX: init,
		maxY: init,
		reader: reader,
	}
	b.setVal(15, 15, 1)  // empty
	intcode := compute.NewIntcode(buf, b)
	if _, err := intcode.Run(); err != nil {
		fmt.Printf("compute.Run: %v\n", err)
		return
	}
	var blocks int
	for _, row := range b.board {
		for _, elem := range row {
			if elem == 2 {
				blocks++
			}
		}
	}
	fmt.Printf("blocks: %d\n", blocks)
}

type board struct {
	posX, posY int
	minX, minY int
	maxX, maxY int

	// 1 - north
	// 2 - south
	// 3 - west
	// 4 - east
	lastMove int64
	
	// 0 - unknown
	// 1 - empty
	// 2 - wall
	// 3 - oxygen
	board [][]int

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
			b.print()
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

func (b board) print() {
	for y := b.minY; y <= b.maxY; y++ {
		for x := b.minX; x <= b.maxX; x++ {
			if y >= len(b.board[x]) {
				fmt.Printf(" ")
				continue
			}
			if x == b.posX && y == b.posY {
				if b.board[x][y] == 3 {
					fmt.Printf("*") // oxygen here
				} else {
					fmt.Printf("+")
				}
				continue
			}
			switch b.board[x][y] {
			case 0:
				fmt.Printf(" ")
			case 1:
				fmt.Printf(".")
			case 2:
				fmt.Printf("X")
			case 3:
				fmt.Printf("O")
			}
		}
		fmt.Printf("\n")
	}
}

func (b *board) setVal(posX, posY, val int) {
	for posX >= len(b.board) {
		b.board = append(b.board, nil)
	}
	for posY >= len(b.board[posX]) {
		b.board[posX] = append(b.board[posX], 0)
	}
	b.board[posX][posY] = val
	if posX < b.minX {
		b.minX = posX
	}
	if posX > b.maxX {
		b.maxX = posX
	}
	if posY < b.minY {
		b.minY = posY
	}
	if posY > b.maxY {
		b.maxY = posY
	}
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
	if b.posX < b.minX {
		b.minX = b.posX
	}
	if b.posX > b.maxX {
		b.maxX = b.posX
	}
	if b.posY < b.minY {
		b.minY = b.posY
	}
	if b.posY > b.maxY {
		b.maxY = b.posY
	}
	
	return nil
}
