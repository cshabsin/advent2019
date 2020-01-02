package main

import (
	"errors"
	"fmt"
	"io"
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
	board := [][]int{nil}
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
			lastRow := len(board) - 1
			if val == 10 {
				board = append(board, nil)
				fmt.Printf("\n")
			} else {
				var bVal int
				if val == 35 { // pound
					bVal = 1
				} else if val == 46 { // period
					bVal = 0
				} else if val == 94 { // caret
					bVal = 2
				} else if val == 60 { // less than
					bVal = 2
				} else if val == 62 { // greater than
					bVal = 2
				} else if val == 118 { // v
					bVal = 2
				}
				board[lastRow] = append(board[lastRow], bVal)
				fmt.Printf("%c", []rune{'.', '#', '*', '?'}[bVal])
			}
		}
	}()
	intcode := compute.NewIntcode(buf, bIO)
	if _, err := intcode.Run(); err != nil {
		fmt.Printf("compute.Run: %v\n", err)
		return
	}
	var total int
	for y, row := range board {
		for x, val := range row {
			if val != 1 {
				continue
			}
			if y == 0 || y == len(board)-1 || x == 0 || x == len(row)-1 {
				// first/last row can't have an intersection
				continue
			}
			if row[x-1] != 1 || row[x+1] != 1 {
				continue
			}
			if board[y-1][x] != 1 || board[y+1][x] != 1 {
				continue
			}
			total += x*y
		}
	}
	fmt.Printf("sum of alignment params: %d\n", total)
}
