package main

import (
	"bufio"
	"errors"
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

	droid := maze.MakeDroid()
	aIO, bIO := compute.NewChanIO()
	//	go ManualRobot(droid, aIO)
	go MappingRobot(droid, aIO)

	intcode := compute.NewIntcode(buf, bIO)
	if _, err := intcode.Run(); err != nil {
		fmt.Printf("compute.Run: %v\n", err)
		return
	}

}

func MappingRobot(d *maze.Droid, io *compute.ChanIO) {
	err := Map(d, io)
	if err != nil {
		fmt.Printf("Map: %v\n", err)
	}
	d.Print()
	io.Close()
}

func Map(d *maze.Droid, io *compute.ChanIO) error {
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
			success, err := d.ProcessMove(dir, io, print)
			if err != nil {
				return err
			}
			if success {
				if err := d.ExpectMove(maze.Opposite(dir), io, print); err != nil {
					return err
				}
				successes = append(successes, dir)
			}
		}
		if len(successes) == 0 {
			for i := len(undoMoves)-1; i>=0; i-- {
				if err := d.ExpectMove(undoMoves[i], io, print); err != nil {
					return err
				}
			}
			return nil
		}
		if len(successes) == 1 {
			if err := d.ExpectMove(successes[0], io, print); err != nil {
				return err
			}
			undoMoves = append(undoMoves, maze.Opposite(successes[0]))
			continue
		}
		for _, nextDir := range successes {
			if err := d.ExpectMove(nextDir, io, print); err != nil {
				return err
			}
			if err := Map(d, io); err != nil {
				return err
			}
			if err := d.ExpectMove(maze.Opposite(nextDir), io, print); err != nil {
				return err
			}
		}
	}
}

func ManualRobot(d *maze.Droid, io *compute.ChanIO) {
	defer io.Close()
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
		_, err = d.ProcessMove(int(move), io, true)
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
