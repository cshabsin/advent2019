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
	b := &board{}
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

	b = &board{}
	buf[0] = 2
	intcode = compute.NewIntcode(buf, b)
	if _, err := intcode.Run(); err != nil {
		fmt.Printf("compute.Run: %v\n", err)
		return
	}
}

type board struct {
	// 0 - posX next
	// 1 - posY next
	// 2 - element next
	phase int

	posX, posY int
	paddleX, ballX int
	
	// 0 - empty
	// 1 - wall
	// 2 - block
	// 3 - horiz paddle
	// 4 - ball
	board [][]int
}

func (b board) Read() (int64, error) {
	var v int64
	if b.paddleX < b.ballX {
		v = 1
	} else if b.paddleX > b.ballX {
		v = -1
	}
	fmt.Printf("input: %d\n", v)
	return v, nil
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
		b.phase = 1
		b.posX = int(val)
	case 1:
		b.phase = 2
		b.posY = int(val)
	case 2:
		b.phase = 0
		if b.posX == -1 {
			if b.posY != 0 {
				return fmt.Errorf("invalid pos: %d, %d", b.posX, b.posY)
			}
			fmt.Printf("score: %d\n", val)
			return nil
		}
		switch val {
		case 3:
			b.paddleX = b.posX
			fmt.Printf("paddle at pos: %d, %d\n", b.posX, b.posY)
		case 4:
			b.ballX = b.posX
			fmt.Printf("ball at pos: %d, %d\n", b.posX, b.posY)
		}
		b.setVal(b.posX, b.posY, int(val))
	}
	return nil
}
