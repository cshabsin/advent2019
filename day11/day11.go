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
	pb := &paintBot{
		posX: 10000,
		posY: 10000,
	}
	intcode := compute.NewIntcode(buf, pb)
	if _, err := intcode.Run(); err != nil {
		fmt.Printf("compute.Run: %v\n", err)
		return
	}
	//	fmt.Printf("out: %v\n", pb.out)
	//      fmt.Printf("board: %v\n", pb.board)
	set := 0
	for _, row := range pb.painted {
		for _, v := range row {
			if v {
				set += 1
			}
		}
	}
	fmt.Printf("set: %d\n", set)
}

type paintBot struct {
	posX, posY int
	dir        int // 0, 1, 2, 3 for north, south, east, west?

	phase   bool // false = next write paints, true = next write turns
	board   [][]bool
	painted [][]bool

	out []int64
}

func (p paintBot) Read() (int64, error) {
	if p.posX >= len(p.board) {
		return 0, nil
	}
	if p.posY >= len(p.board[p.posX]) {
		return 0, nil
	}
	if p.board[p.posX][p.posY] {
		return 1, nil
	}
	return 0, nil
}

func (p *paintBot) Write(val int64) error {
	p.out = append(p.out, val)
	if p.phase {
		switch val {
		case 0:
			p.dir = p.dir - 1
			if p.dir < 0 {
				p.dir += 4
			}
		case 1:
			p.dir = (p.dir + 1) % 4
		default:
			return fmt.Errorf("bad val during turn phase: %d", val)
		}
		switch p.dir {
		case 0: // north
			p.posY += 1

		case 1: // east
			p.posX += 1
		case 2: // south
			p.posY -= 1
		case 3: // west
			p.posX -= 1
		default:
			return fmt.Errorf("invalid dir %d", p.dir)
		}
	} else {
		for p.posX >= len(p.board) {
			p.board = append(p.board, []bool{})
		}
		for p.posY >= len(p.board[p.posX]) {
			p.board[p.posX] = append(p.board[p.posX], false)
		}
		for p.posX >= len(p.painted) {
			p.painted = append(p.painted, []bool{})
		}
		for p.posY >= len(p.painted[p.posX]) {
			p.painted[p.posX] = append(p.painted[p.posX], false)
		}
		p.painted[p.posX][p.posY] = true
		if val == 0 {
			p.board[p.posX][p.posY] = false
		} else if val == 1 {
			p.board[p.posX][p.posY] = true
		} else {
			return fmt.Errorf("bad val during paint phase: %d", val)
		}
	}
	p.phase = !p.phase
	return nil
}
