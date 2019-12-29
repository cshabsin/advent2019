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
	board := &board{}
	intcode := compute.NewIntcode(buf, board)
	if _, err := intcode.Run(); err != nil {
		fmt.Printf("compute.Run: %v\n", err)
		return
	}
	var blocks int
	for _, row := range board.board {
		for _, elem := range row {
			if elem == 2 {
				blocks++
			}
		}
	}
	fmt.Printf("blocks: %d\n", blocks)
}

type board struct {
	// 0 - posX next
	// 1 - posY next
	// 2 - element next
	phase int

	posX, posY int
	
	// 0 - empty
	// 1 - wall
	// 2 - block
	// 3 - horiz paddle
	// 4 - ball
	board [][]int
}

func (b board) Read() (int64, error) {
	return 0, nil
}

func (b *board) setVal(posX, posY int, val int) {
	for posX >= len(b.board) {
		b.board = append(b.board, nil)
	}
	for posY >= len(b.board[posX]) {
		b.board[posX] = append(b.board[posX], 0)
	}
	b.board[posX][posY] = val
}

func (b *board) Write(val int64) error {
	switch b.phase {
	case 0:
		b.posX = int(val)
		b.phase = 1
	case 1:
		b.posY = int(val)
		b.phase = 2
	case 2:
		b.setVal(b.posX, b.posY, int(val))
		b.phase = 0
	}
	return nil
}
