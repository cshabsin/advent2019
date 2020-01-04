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

	cmdChan    chan string
	chanActive bool
	hasCmd     bool
	cmd        string
	cmdIndex   int

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
		if a.cmdChan == nil {
			a.cmdChan = make(chan string, 1)
		}
		if a.chanActive {
			cmd := <-a.cmdChan
			fmt.Printf("> %s\n", cmd)
			a.SetCmd(cmd)
		} else {
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
	fmt.Printf("> %s\n", cmds[0])
	a.SetCmd(cmds[0])
	a.chanActive = true
	go func() {
		for _, cmd := range cmds[1:] {
			a.cmdChan <- cmd
		}
		a.chanActive = false
	}()
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
	a.chanActive = true
	go func() {
		for _ = range items {
		}
		a.chanActive = false
	}()
}

func (a *AsciiIO) ReadFromUser() error {
	fmt.Print("> ")
	move, err := a.reader.ReadString('\n')
	if err != nil {
		return err
	}
	move = strings.TrimSpace(move)
	if move == "state" {
		fmt.Printf("state: %s\n", a.state)
	} else if move == "n" {
		a.SetCmd("north")
	} else if move == "s" {
		a.SetCmd("south")
	} else if move == "e" {
		a.SetCmd("east")
	} else if move == "w" {
		a.SetCmd("west")
	} else if move == "dropall" {
		var cmds []string
		for _, i := range items {
			cmds = append(cmds, "drop "+i)
		}
		a.StreamCmds(cmds)
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
	} else if move == "search" {
		a.Search()
	} else {
		a.SetCmd(move)
	}
	return nil
}

func (a *AsciiIO) Close() {
}
