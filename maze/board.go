package maze

import "fmt"

type Board struct {
	// offsets, to allow for negative logical coords
	offsetX, offsetY int

	// 0 - unknown
	// 1 - empty
	// 2 - wall
	// 3 - oxygen
	board [][]int
}

func NewBoard() *Board {
	b := Board{
		board: [][]int{},
	}
	b.SetVal(0, 0, 1)
	return &b
}

func (b *Board) padTop(yPad int) {
	newBoard := make([][]int, 0, len(b.board)+yPad)
	for i := 0; i < yPad; i++ {
		newBoard = append(newBoard, nil)
	}
	for _, row := range b.board {
		newBoard = append(newBoard, row)
	}
	b.offsetY += yPad
	b.board = newBoard
}

func (b *Board) padLeft(xPad int) {
	newBoard := make([][]int, 0, len(b.board))
	for _, row := range b.board {
		var newRow []int
		if row != nil {
			newRow = make([]int, 0, len(row)+xPad)
			for i := 0; i < xPad; i++ {
				newRow = append(newRow, 0)
			}
			newRow = append(newRow, row...)
		}
		newBoard = append(newBoard, newRow)
	}
	b.offsetX += xPad
	b.board = newBoard
}

func (b *Board) SetVal(posX, posY, val int) {
	x := posX + b.offsetX
	if x < 0 {
		b.padLeft(-x)
		// increase offsetX, rewrite arrays
		x = posX + b.offsetX
	}
	y := posY + b.offsetY
	if y < 0 {
		b.padTop(-y)
		// increase offsetY, rewrite arrays
		y = posY + b.offsetY
	}
	for y >= len(b.board) {
		b.board = append(b.board, nil)
	}
	for x >= len(b.board[y]) {
		b.board[y] = append(b.board[y], 0)
	}
	b.board[y][x] = val
}

func (b Board) GetVal(posX, posY int) int {
	y := posY + b.offsetY
	if y < 0 || y >= len(b.board) {
		return 0
	}
	x := posX + b.offsetX
	if x < 0 || x >= len(b.board[y]) {
		return 0
	}
	return b.board[y][x]
}

// Edges returns minX, minY, maxX, maxY
func (b Board) Edges() (int, int, int, int) {
	width := 0
	for _, row := range b.board {
		if len(row) > width {
			width = len(row)
		}
	}
	return -b.offsetX, -b.offsetY, len(b.board) - b.offsetX, width - b.offsetY
}

func (b Board) Print() {
	b.print(0, 0, false)
}

func (b Board) print(droidX, droidY int, include bool) {
	for y := 0; y < len(b.board); y++ {
		for x := 0; x < len(b.board[y]); x++ {
			if include && droidX+b.offsetX == x && droidY+b.offsetY == y {
				if b.board[y][x] == 3 {
					fmt.Printf("*")
				} else {
					fmt.Printf("+")
				}
			} else {
				switch b.board[y][x] {
				case 0:
					fmt.Printf(" ")
				case 1:
					fmt.Printf(".")
				case 2:
					fmt.Printf("X")
				case 3:
					fmt.Printf("O")
				}
			}
		}
		fmt.Printf("\n")
	}
}

func (b Board) Board() [][]int {
	return b.board
}
