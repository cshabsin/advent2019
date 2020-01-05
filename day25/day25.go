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
	fmt.Printf("g(0): %v, g(1): %v, g(2): %v\n", Gray(0), Gray(1), Gray(2))
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
		echo:    true,
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

	echo bool

	cmdChan  chan string
	eof      chan bool
	hasCmd   bool
	cmd      string
	cmdIndex int

	readTrigger bool
	readSem     chan bool

	state string
	room  string
}

func (a *AsciiIO) Write(val int64) error {
	if val < 0 || val > 255 {
		return fmt.Errorf("invalid Write: %d", val)
	}
	if val == 10 {
		if strings.Contains(a.lastLine, "lighter") {
			a.state = "too heavy"
		} else if strings.Contains(a.lastLine, "heavier") {
			a.state = "too light"
		}
		if strings.HasPrefix(a.lastLine, "== ") {
			a.room = a.lastLine[3 : len(a.lastLine)-3]
		}
		a.lastLine = ""
	} else {
		a.lastLine += string(rune(val))
	}
	if a.echo {
		fmt.Printf("%c", rune(val))
	}
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
	"astronaut ice cream", // 0 -> 1
	"mouse",               // 1 -> 2
	"ornament",            // 2 -> 4
	"easter egg",          // 3 -> 8
	"hypercube",           // 4 -> 16
	"prime number",        // 5 -> 32
	"wreath",              // 6 -> 64
	"mug",                 // 7 -> 128
}

// returns the list of bits to permute for an n-bit gray code.
func Gray(n int) []int {
	if n == 0 {
		return []int{0}
	}
	g := Gray(n - 1)
	l := len(g)
	g = append(g, n)
	for i := 0; i < l; i++ {
		g = append(g, g[l-1-i])
	}
	return g
}

func (a *AsciiIO) Search() {
	go func() {
		//		a.echo = false
		found := make([]int, 256)
		inv := 255
		fmt.Printf("Gray(7): %v\n", Gray(7))
		for _, bit := range Gray(7) {
			a.state = ""
			if inv&(1<<bit) != 0 {
				a.SendCmd("drop "+items[bit], true)
			} else {
				a.SendCmd("take "+items[bit], true)
			}
			inv ^= 1 << bit
			a.SendCmd("north", true)
			a.Wait()
			if a.state == "" {
				found[inv] = 3
				fmt.Printf("found it: %d\n", inv)
				fmt.Print(a.lastMsg)
				a.echo = true
				a.SendCmd("inv", true)
				return
			} else if a.state == "too light" {
				found[inv] = 2
			} else if a.state == "too heavy" {
				found[inv] = 1
			}
		}
		a.echo = true
		fmt.Printf("found: %v\n", found)
	}()
}

func (a *AsciiIO) Wait() {
	if a.readTrigger {
		<-a.readSem
		a.readTrigger = false
	}
}

func (a *AsciiIO) SendCmd(cmd string, print bool) {
	a.Wait()
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
		} else if move == "echo" {
			a.echo = true
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
