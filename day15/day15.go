package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

	aIO, bIO := compute.NewChanIO()
	droid := maze.MakeDroid(aIO)
	go MappingRobot(droid)

	intcode := compute.NewIntcode(buf, bIO)
	if _, err := intcode.Run(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Printf("compute.Run: %v\n", err)
		return
	}

	if !droid.OxygenFound() {
		fmt.Printf("oxygen not found\n")
		return
	}

	board := droid.Board()
	minX, minY, maxX, maxY := board.Edges()
	var steps int
	for {
		var found bool
		for y := minY; y <= maxY; y++ {
			for x := minX; x <= maxX; x++ {
				if board.GetVal(x, y) == 3 {
					for _, dir := range maze.Dirs {
						if board.GetVal(x+dir.DX, y+dir.DY) == 1 {
							board.SetVal(x+dir.DX, y+dir.DY, 4)
							found = true
						}
					}
				}
			}
		}
		for y := minY; y <= maxY; y++ {
			for x := minX; x <= maxX; x++ {
				if board.GetVal(x, y) == 4 {
					board.SetVal(x, y, 3)
				}
			}
		}
		if !found {
			break
		}
		steps++
	}
	board.Print()
	fmt.Printf("steps: %d\n", steps)
}

func MappingRobot(d *maze.Droid) {
	err := Map(d)
	if err != nil {
		fmt.Printf("Map: %v\n", err)
	}
	d.Print()
	d.CloseIO()
}

func Map(d *maze.Droid) error {
	print := false
	var undoMoves []int
	for {
		if print {
			fmt.Print("\n")
		}
		var successes []int
		for dir := 1; dir <= 4; dir++ {
			if d.LookDir(dir) != 0 {
				continue // already known
			}
			success, err := d.ProcessMove(dir, print)
			if err != nil {
				return err
			}
			if success {
				if err := d.ExpectMove(maze.Opposite(dir), print); err != nil {
					return err
				}
				successes = append(successes, dir)
			}
		}
		if len(successes) == 0 {
			for i := len(undoMoves) - 1; i >= 0; i-- {
				if err := d.ExpectMove(undoMoves[i], print); err != nil {
					return err
				}
			}
			return nil
		}
		if len(successes) == 1 {
			if err := d.ExpectMove(successes[0], print); err != nil {
				return err
			}
			undoMoves = append(undoMoves, maze.Opposite(successes[0]))
			continue
		}
		for _, nextDir := range successes {
			if err := d.ExpectMove(nextDir, print); err != nil {
				return err
			}
			if err := Map(d); err != nil {
				return err
			}
			if err := d.ExpectMove(maze.Opposite(nextDir), print); err != nil {
				return err
			}
		}
	}
}

func ManualRobot(d *maze.Droid) {
	defer d.CloseIO()
	reader := bufio.NewReader(os.Stdin)
	for {
		move, err := ReadMove(reader)
		var m *ManualInputError
		if errors.As(err, &m) {
			if m.input == "p" {
				d.Print()
			} else {
				fmt.Printf("invalid input %q\n", m.input)
			}
			continue
		}
		if err != nil {
			fmt.Printf("ReadMove: %v\n", err)
			return
		}
		_, err = d.ProcessMove(int(move), true)
		if err != nil {
			fmt.Printf("droid.ProcessMove: %v\n", err)
			return
		}
	}
}

type ManualInputError struct {
	input string
}

func (m ManualInputError) Error() string {
	return fmt.Sprintf("manual input %q", m.input)
}

func ReadMove(reader *bufio.Reader) (int64, error) {
	fmt.Printf("input: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	text = strings.TrimSpace(text)
	if text == "n" {
		return 1, nil
	} else if text == "s" {
		return 2, nil
	} else if text == "w" {
		return 3, nil
	} else if text == "e" {
		return 4, nil
	} else {
		val, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return 0, &ManualInputError{text}
		}
		return val, nil
	}
}
