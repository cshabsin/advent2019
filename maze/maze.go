package maze

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
	return &Board{
		board: [][]int{},
	}
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
	b.board[x][y] = val
}
