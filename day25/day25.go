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
	"time"

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

	aIO, bIO := compute.NewChanIO()
	aIO.NoTimeout = true
	bIO.NoTimeout = true

	var state, line string

	go func() {
		for {
			val, err := aIO.Read()
			if errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				fmt.Printf("read: %v\n", err)
				break
			} else if val < 256 {
				fmt.Printf("%c", rune(val))
				if val == 10 {
					if strings.Contains(line, "lighter") {
						state = "lighter"
					} else if strings.Contains(line, "heavier") {
						state = "heavier"
					}
					line = ""
				} else {
					line += string(rune(val))
				}
			} else {
				fmt.Printf("Value: %d\n", val)
			}
		}
	}()

	intcode := compute.NewIntcode(buf, bIO)
	done := false
	go func() {
		if _, err := intcode.Run(); err != nil {
			fmt.Printf("Run: %v", err)
		}
		done = true
	}()
	defer aIO.Close()

	cmds := []string{
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
	for _, cmd := range cmds {
		for {
			time.Sleep(50 * time.Millisecond)
			if aIO.Idle() {
				break
			}
		}
		fmt.Printf("> %s\n", cmd)
		for _, r := range cmd + "\n" {
			if err := aIO.Write(int64(r)); err != nil {
				fmt.Printf("Write: %v\n", err)
				break
			}
		}
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		if done {
			break
		}
		for !aIO.Idle() {
			time.Sleep(50 * time.Millisecond)
		}
		fmt.Print("> ")
		move, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("ReadString: %v\n", err)
			break
		}
		if strings.TrimSpace(move) == "state" {
			fmt.Printf("State: %s\n", state)
		}
		for _, r := range move {
			if err := aIO.Write(int64(r)); err != nil {
				fmt.Printf("Write: %v\n", err)
				break
			}
		}
	}
}

func doRead(aIO *compute.ChanIO) {
	val, err := aIO.Read()
	if errors.Is(err, io.EOF) {
		return
	} else if err != nil {
		fmt.Printf("read: %v\n", err)
	} else if val < 256 {
		fmt.Printf("%c", rune(val))
	} else {
		fmt.Printf("Value: %d\n", val)
	}
}

type AsciiIO struct {
	reader  *bufio.Reader
	lastMsg string

	cmds     []string
	cmdIndex int
}

func (a *AsciiIO) Write(val int64) error {
	if val < 0 || val > 255 {
		return fmt.Errorf("invalid Write: %d", val)
	}
	fmt.Printf("%c", rune(val))
	a.lastMsg += string(rune(val))
	return nil
}

func (a *AsciiIO) Read(val int64) (int64, error) {
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
	a.cmdIndex = 0
	fmt.Print("> ")
	move, err := a.reader.ReadString('\n')
	if err != nil {
		return err
	}
	move = strings.TrimSpace(move)
	if move == "go" {
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
	} else {
		a.cmds = []string{move}
	}
	return nil
}
