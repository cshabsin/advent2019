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
	reader   *bufio.Reader
	lastMsg  string
	lastLine string

	cmdChan  chan string // buffer of 1
	hasCmd   bool
	cmd      string
	cmdIndex int

	state string
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
		if a.cmdChan == nil {
			a.cmdChan = make(chan string, 1)
		}
		select {
		case cmd := <- a.cmdChan:
			fmt.Printf("> %s\n", cmd)
			a.SetCmd(cmd)
		default:
			if err := a.ReadFromUser(); err != nil {
				return 0, err
			}
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

func (a *AsciiIO) SetCmd(cmd string) {
	a.cmdIndex = 0
	a.cmd = cmd
	a.hasCmd = true
}

func (a *AsciiIO) StreamCmds(cmds []string) {
	ready := make(chan bool)
	go func() {
		sent := false
		for _, cmd := range cmds {
			a.cmdChan <- cmd
			if !sent {
				sent = true
				ready <- true
			}
		}
	}()
	<-ready
}

func (a *AsciiIO) ReadFromUser() error {
	for {
		fmt.Print("> ")
		move, err := a.reader.ReadString('\n')
		if err != nil {
			return err
		}
		move = strings.TrimSpace(move)
		if move == "state" {
			fmt.Printf("state: %s\n", a.state)
			return nil
		} else if move == "n" {
			a.SetCmd("north")
			return nil
		} else if move == "s" {
			a.SetCmd("south")
			return nil
		} else if move == "e" {
			a.SetCmd("east")
			return nil
		} else if move == "w" {
			a.SetCmd("west")
			return nil
		} else if move == "dropall" {
			a.StreamCmds([]string{
				"drop astronaut ice cream",
				"drop mouse",
				"drop ornament",
				"drop easter egg",
				"drop hypercube",
				"drop prime number",
				"drop wreath",
				"drop mug",
			})
			return nil
		} else if move == "go" {
			a.StreamCmds([]string{
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
			})
			return nil
		} else {
			a.SetCmd(move)
			return nil
		}
	}
}

func (a *AsciiIO) Close() {
}
