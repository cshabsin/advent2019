package maze

import (
	"fmt"
	"time"

	"github.com/cshabsin/advent2019/compute"
)

type Droid struct {
	posX, posY int
	board      *Board
}

func MakeDroid() *Droid {
	return &Droid{
		board: NewBoard(),
	}
}

func (d *Droid) ProcessMove(move int64, io *compute.ChanIO) error {
	select {
	case io.Input <- move:
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout injecting into Input")
	}
	val := <-io.Output
	switch move {
	case 1: // north
		d.Move(0, -1, val, true)
	case 2: // south
		d.Move(0, 1, val, true)
	case 3: // west
		d.Move(-1, 0, val, true)
	case 4: // east
		d.Move(1, 0, val, true)
	}
	return nil
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
