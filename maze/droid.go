package maze

import "fmt"

type Droid struct {
	posX, posY int
	board      *Board
}

func MakeDroid() *Droid {
	return &Droid{
		board: NewBoard(),
	}
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
