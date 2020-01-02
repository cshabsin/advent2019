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

func main() {
	content, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	buf, err := compute.ParseFile(content)
	if err != nil {
		log.Fatal(err)
	}
	MoveBot(buf)
}

func Map(buf []int64) {
	s := scaffold{}
	aIO, bIO := compute.NewChanIO()
	go func() {
		for {
			val, err := aIO.Read()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				fmt.Printf("aIO.Read: %v", err)
				return
			}
			if err := s.populate(val); err != nil {
				fmt.Printf("s.populate: %v", err)
			}
			fmt.Printf("%c", rune(val))
		}
	}()
	intcode := compute.NewIntcode(buf, bIO)
	if _, err := intcode.Run(); err != nil {
		fmt.Printf("compute.Run: %v\n", err)
		return
	}

	fmt.Printf("bot at %d, %d (direction %s)\n", s.posX, s.posY, s.dir)
	path := s.walk()
	fmt.Printf("path: %v\n", path)
	a := []string{"R", "12", "R", "4", "L", "6", "L", "8", "L", "8"}
	path = substitute(path, a, []string{"A"})
	fmt.Printf("path after A: %v\n", path)

	b := []string{"L", "12", "R", "4", "R", "4"}
	path = substitute(path, b, []string{"B"})
	fmt.Printf("path after B: %v\n", path)

	c := []string{"R", "12", "R", "4", "L", "12"}
	path = substitute(path, c, []string{"C"})
	fmt.Printf("path after C: %v\n", path)
}

func MoveBot(buf []int64) {
	a := []string{"R", "12", "R", "4", "L", "6", "L", "8", "L", "8"}
	b := []string{"L", "12", "R", "4", "R", "4"}
	c := []string{"R", "12", "R", "4", "L", "12"}
	path := []string{"B", "C", "C", "A", "A", "B", "B", "C", "C", "A"}
	buf[0] = 2 // move bot mode
	aIO, bIO := compute.NewChanIO()
	fin := make(chan bool)
	go func() {
		for {
			val, err := aIO.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				fmt.Printf("read: %v\n", err)
				return
			}
			if val < 256 {
				fmt.Printf("%c", rune(val))
			} else {
				fmt.Printf("read: %d\n", val)
			}
		}
		fin <- true
	}()
	go func() {
		defer aIO.Close()
		if err := output(aIO, path); err != nil {
			fmt.Printf("output path: %v\n", err)
			return
		}
		if err := output(aIO, a); err != nil {
			fmt.Printf("output A: %v\n", err)
			return
		}
		if err := output(aIO, b); err != nil {
			fmt.Printf("output B: %v\n", err)
			return
		}
		if err := output(aIO, c); err != nil {
			fmt.Printf("output C: %v\n", err)
			return
		}
		if err := output(aIO, []string{"n"}); err != nil {
			fmt.Printf("output n: %v\n", err)
			return
		}
	}()
	intcode := compute.NewIntcode(buf, bIO)
	if _, err := intcode.Run(); err != nil {
		fmt.Printf("compute.Run: %v\n", err)
		return
	}
	<-fin
}

func output(io compute.IO, s []string) error {
	sout := strings.Join(s, ",") + "\n"
	fmt.Printf("%s", sout)
	for _, r := range sout {
		if err := io.Write(int64(r)); err != nil {
			return err
		}
	}
	return nil
}

func printTotal(s *scaffold) {
	var total int
	for y, row := range s.board {
		for x, val := range row {
			if val != 1 {
				continue
			}
			if y == 0 || y == len(s.board)-1 || x == 0 || x == len(row)-1 {
				// first/last row can't have an intersection
				continue
			}
			if row[x-1] != 1 || row[x+1] != 1 {
				continue
			}
			if s.board[y-1][x] != 1 || s.board[y+1][x] != 1 {
				continue
			}
			total += x * y
		}
	}
	fmt.Printf("sum of alignment params: %d\n", total)
}

var dirs = []struct{ dx, dy int }{
	{1, 0},  // east
	{0, -1}, // north
	{-1, 0}, // west
	{0, 1},  // south
}

type direction int

func (d direction) turn(left int) direction {
	return direction((int(d) + left) % 4)
}

func (d direction) left() direction {
	return d.turn(1)
}

func (d direction) opposite() direction {
	return d.turn(2)
}

func (d direction) right() direction {
	return d.turn(3)
}

func (d direction) dxy() (int, int) {
	return dirs[int(d)].dx, dirs[int(d)].dy
}

func (d direction) String() string {
	switch d {
	case east:
		return "east"
	case north:
		return "north"
	case west:
		return "west"
	case south:
		return "south"
	default:
		return fmt.Sprintf("bad(%d)", int(d))
	}
}

const (
	east  = direction(0)
	north = direction(1)
	west  = direction(2)
	south = direction(3)
)

type scaffold struct {
	board [][]int

	found      bool
	posX, posY int
	dir        direction
}

func (s *scaffold) populateBot(d direction) error {
	y := len(s.board) - 1
	x := len(s.board[y])
	if s.found {
		return fmt.Errorf("found second bot at %d, %d (first at %d, %d)", x, y, s.posX, s.posY)
	}
	s.found = true
	s.posX = x
	s.posY = y
	s.dir = d
	return nil
}

func (s *scaffold) populate(val int64) error {
	if s.board == nil {
		s.board = [][]int{nil}
	}
	lastRow := len(s.board) - 1
	if val == 10 {
		s.board = append(s.board, nil)
	} else {
		var bVal int
		if val == 35 { // pound
			bVal = 1
		} else if val == 46 { // period
			bVal = 0
		} else if val == 94 { // caret
			if err := s.populateBot(north); err != nil {
				return err
			}
			bVal = 2
		} else if val == 60 { // less than
			if err := s.populateBot(west); err != nil {
				return err
			}
			bVal = 2
		} else if val == 62 { // greater than
			if err := s.populateBot(east); err != nil {
				return err
			}
			bVal = 2
		} else if val == 118 { // v
			if err := s.populateBot(south); err != nil {
				return err
			}
			bVal = 2
		} else {
			return fmt.Errorf("invalid value received: %d", val)
		}
		s.board[lastRow] = append(s.board[lastRow], bVal)
	}
	return nil
}

func (s scaffold) canMove(dir direction) int {
	move := 0
	posX := s.posX
	posY := s.posY
	dx, dy := dir.dxy()
	for {
		posX += dx
		posY += dy
		if posX < 0 || posY < 0 || posY >= len(s.board) || posX >= len(s.board[posY]) || s.board[posY][posX] == 0 {
			return move
		}
		move++
	}
}

// Returns a list of moves (L, R, or number)
func (s *scaffold) walk() []string {
	var moves []string
	for {
		if move := s.canMove(s.dir); move != 0 {
			moves = append(moves, fmt.Sprintf("%d", move))
			dx, dy := s.dir.dxy()
			s.posX += move * dx
			s.posY += move * dy
		} else {
			if s.canMove(s.dir.right()) != 0 {
				moves = append(moves, "R")
				s.dir = s.dir.right()
			} else if s.canMove(s.dir.left()) != 0 {
				moves = append(moves, "L")
				s.dir = s.dir.left()
			} else {
				return moves
			}
		}

	}
}

func equals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func substitute(path, subpath, repl []string) []string {
	var newPath []string
	for i := 0; i < len(path); {
		if i+len(subpath) <= len(path) && equals(path[i:i+len(subpath)], subpath) {
			newPath = append(newPath, repl...)
			i += len(subpath)
		} else {
			newPath = append(newPath, path[i])
			i++
		}
	}
	return newPath
}
