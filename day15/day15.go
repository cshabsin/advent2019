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
	io := compute.NewChanIO()
	go func(d *maze.Droid, io *compute.ChanIO) {
		defer io.Close()
		reader := bufio.NewReader(os.Stdin)
		for {
			move, err := ReadMove(reader)
			if err == printBoard {
				d.Print()
				continue
			}
			if err != nil {
				fmt.Printf("ReadMove: %v\n", err)
				return
			}
			io.Input <- move
			val := <-io.Output
			if err != nil {
				fmt.Printf("io.Read: %v\n", err)
				return
			}
			switch move {
			case 1: // north
				droid.Move(0, -1, val, true)
			case 2: // south
				droid.Move(0, 1, val, true)
			case 3: // west
				droid.Move(-1, 0, val, true)
			case 4: // east
				droid.Move(1, 0, val, true)
			}
		}
	}(droid, io)
	intcode := compute.NewIntcode(buf, io)
	if _, err := intcode.Run(); err != nil {
		fmt.Printf("compute.Run: %v\n", err)
		return
	}
}

var printBoard = errors.New("print board requested")

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
	} else if text == "p" {
		return 0, printBoard
	} else {
		return strconv.ParseInt(strings.TrimSpace(text), 10, 64)
	}
}
