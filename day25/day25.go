package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	io := &AsciiIO{reader: reader}
	intcode := compute.NewIntcode(buf, io)
	if _, err := intcode.Run(); err != nil {
		fmt.Printf("Run: %v\n", err)
	}
}

type AsciiIO struct {
	reader  *bufio.Reader
	lastMsg string

	cmds     []string
	cmdIndex int

	state string
}

func (a *AsciiIO) Write(val int64) error {
	if val < 0 || val > 255 {
		return fmt.Errorf("invalid Write: %d", val)
	}
	fmt.Printf("%c", rune(val))
	a.lastMsg += string(rune(val))
	return nil
}

func (a *AsciiIO) Read() (int64, error) {
	if len(a.cmds) == 0 {
		if err := a.ReadFromUser(); err != nil {
			return 0, err
		}
	}
	if a.cmdIndex == len(a.cmds[0]) {
		a.cmds = a.cmds[1:]
		a.cmdIndex = 0
		return 10, nil
	}
	rc := int64(a.cmds[0][a.cmdIndex])
	a.cmdIndex++
	return rc, nil
}

func (a *AsciiIO) ReadFromUser() error {
	if strings.Contains(a.lastMsg, "lighter") {
		a.state = "lighter"
	} else if strings.Contains(a.lastMsg, "heavier") {
		a.state = "heavier"
	}

	a.cmdIndex = 0
	for {
		fmt.Print("> ")
		move, err := a.reader.ReadString('\n')
		if err != nil {
			return err
		}
		move = strings.TrimSpace(move)
		if move == "state" {
			fmt.Printf("state: %s\n", a.state)
		} else if move == "go" {
			a.cmds = []string{
				"north",
				"take astronaut ice cream",
				"south",
				"west",
				"take mouse",
				"north",
				"take ornament",
				"west",
				"north",
				"take easter egg",
				"east",
				"take hypercube",
				"north",
				"east",
				"take prime number",
				"west",
				"south",
				"west",
				"north",
				"west",
				"north",
				"take wreath",
				"south",
				"east",
				"south",
				"south",
				"west", // don't take mug, it's too heavy
				"west",
				"inv",
			}
			return nil
		} else {
			a.cmds = []string{move}
			return nil
		}
	}
}

func (a *AsciiIO) Close() {
}
