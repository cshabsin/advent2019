package maze

import (
	"fmt"

	"github.com/cshabsin/advent2019/compute"
)

type Droid struct {
	io         *compute.ChanIO
	posX, posY int
	board      *Board
}

func MakeDroid(io *compute.ChanIO) *Droid {
	return &Droid{
		io:    io,
		board: NewBoard(),
	}
}

func Opposite(move int) int {
	switch move {
	case 1:
		return 2
	case 2:
		return 1
	case 3:
		return 4
	case 4:
		return 3
	}
	fmt.Printf("Opposite: unexpected input %d\n", move)
	return 0
}

var dirs = map[int]struct{ dx, dy int }{
	1: {0, -1},
	2: {0, 1},
	3: {-1, 0},
	4: {1, 0}}

func (d Droid) LookDir(move int) int {
	diff := dirs[move]
	posX := d.posX + diff.dx
	posY := d.posY + diff.dy
	return d.board.GetVal(posX, posY)
}

func (d *Droid) ExpectMove(move int, print bool) error {
	success, err := d.ProcessMove(move, print)
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("failure of expected successful move")
	}
	return nil
}

// Return true if move successful
func (d *Droid) ProcessMove(move int, print bool) (bool, error) {
	err := d.io.Write(int64(move))
	if err != nil {
		return false, err
	}
	val, err := d.io.Read()
	if err != nil {
		return false, err
	}
	diff := dirs[move]
	d.Move(diff.dx, diff.dy, val, print)
	return val != 0, nil
}

func (d *Droid) Move(dx, dy int, result int64, print bool) {
	nextPosX := d.posX + dx
	nextPosY := d.posY + dy
	switch result {
	case 0: // wall, no move
		if print {
			fmt.Printf("wall at %d, %d; no move\n", nextPosX, nextPosY)
		}
		d.board.SetVal(nextPosX, nextPosY, 2)
	case 1: // move
		if print {
			fmt.Printf("moved to %d, %d\n", nextPosX, nextPosY)
		}
		d.board.SetVal(nextPosX, nextPosY, 1)
		d.posX = nextPosX
		d.posY = nextPosY
	case 2: // move, oxygen
		if print {
			fmt.Printf("oxygen at %d, %d; moved in\n", nextPosX, nextPosY)
		}
		d.board.SetVal(nextPosX, nextPosY, 3)
		d.posX = nextPosX
		d.posY = nextPosY
	}
}

func (d Droid) Print() {
	d.board.print(d.posX, d.posY, true)
}

func (d *Droid) CloseIO () {
	d.io.Close()
}
