package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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
	io := &AsciiIO{
		reader:  reader,
		cmdChan: make(chan string),
		readSem: make(chan bool),
		eof:     make(chan bool),
	}
	go io.ReadLoop()
	intcode := compute.NewIntcode(buf, io)
	if _, err := intcode.Run(); err != nil {
		fmt.Printf("Run: %v\n", err)
	}
}

type AsciiIO struct {
	reader   *bufio.Reader
	lastMsg  string
	lastLine string

	cmdChan  chan string
	eof      chan bool
	hasCmd   bool
	cmd      string
	cmdIndex int

	readTrigger bool
	readSem  chan bool

	state string
	room  string
}

func (a *AsciiIO) Write(val int64) error {
	// add echo state
	if val < 0 || val > 255 {
		return fmt.Errorf("invalid Write: %d", val)
	}
	if val == 10 {
		if strings.Contains(a.lastLine, "lighter") {
			a.state = "lighter"
		} else if strings.Contains(a.lastLine, "heavier") {
			a.state = "heavier"
		}
		if strings.HasPrefix(a.lastLine, "== ") {
			a.room = a.lastLine[3 : len(a.lastLine)-3]
		}
		a.lastLine = ""
	} else {
		a.lastLine += string(rune(val))
	}
	fmt.Printf("%c", rune(val))
	a.lastMsg += string(rune(val))
	return nil
}

func (a *AsciiIO) Read() (int64, error) {
	for !a.hasCmd {
		a.readSem <- true
		select {
		case cmd := <-a.cmdChan:
			a.cmdIndex = 0
			a.cmd = cmd
			a.hasCmd = true
		case <-a.eof:
			return 0, io.EOF
		}
	}
	if a.cmdIndex == len(a.cmd) {
		a.hasCmd = false
		return 10, nil
	}
	rc := int64(a.cmd[a.cmdIndex])
	a.cmdIndex++
	return rc, nil
}

var items = []string{
	"astronaut ice cream",
	"mouse",
	"ornament",
	"easter egg",
	"hypercube",
	"prime number",
	"wreath",
	"mug",
}

func (a *AsciiIO) Search() {
	go func() {
		for _ = range items {
		}
	}()
}

func (a *AsciiIO) SendCmd(cmd string, print bool) {
	if a.readTrigger {
		<-a.readSem
		a.readTrigger = false
	}
	if print {
		fmt.Printf("> %s\n", cmd)
	}
	a.cmdChan <- cmd
	a.readTrigger = true
}

func (a *AsciiIO) ReadLoop() {
	a.readTrigger = true
	for {
		if a.readTrigger {
			<-a.readSem
			a.readTrigger = false
		}
		fmt.Print("> ")
		move, err := a.reader.ReadString('\n')
		if err != nil {
			a.eof <- true
			if !errors.Is(err, io.EOF) {
				fmt.Printf("ReadString: %v\n", err)
			}
			return
		}
		move = strings.TrimSpace(move)
		if move == "state" {
			fmt.Printf("state: %s\n", a.state)
		} else if move == "n" {
			a.SendCmd("north", false)
		} else if move == "s" {
			a.SendCmd("south", false)
		} else if move == "e" {
			a.SendCmd("east", false)
		} else if move == "w" {
			a.SendCmd("west", false)
		} else if move == "dropall" {
			for _, i := range items {
				a.SendCmd("drop "+i, true)
			}
		} else if move == "go" {
			for _, cmd := range []string{
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
				"west",
				"take mug",
				"west",
				"inv",
			} {
				a.SendCmd(cmd, true)
			}
		} else if move == "search" {
			a.Search()
		} else {
			a.SendCmd(move, false)
		}
	}
}

func (a *AsciiIO) Close() {
}
