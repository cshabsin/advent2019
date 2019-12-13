package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	content, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(content), "\n")
	p0, err := makePath(lines[0])
	if err != nil {
		log.Fatal(err)
	}
	res, steps, err := p0.intersect(lines[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("result: %d\n", res)
	fmt.Printf("steps: %d\n", steps)
}

type path struct {
	// set of columns that are used for each row
	// value of each cell is # of steps to reach cell
	rows map[int]map[int]int
}

func parseEntry(entry string) (int, int, int, error) {
	dx := 0
	dy := 0
	switch entry[0] {
	case 'R':
		dx = 1
	case 'L':
		dx = -1
	case 'U':
		dy = -1
	case 'D':
		dy = 1
	default:
		return 0, 0, 0, fmt.Errorf("invalid dir in entry %s", entry)
	}
	mag, err := strconv.Atoi(entry[1:])
	if err != nil {
		return 0, 0, 0, err
	}
	return dx, dy, mag, nil
}

func makePath(line string) (*path, error) {
	row := 0
	col := 0
	steps := 0
	p := path{rows: map[int]map[int]int{}}
	for _, entry := range strings.Split(line, ",") {
		dx, dy, mag, err := parseEntry(entry)
		if err != nil {
			return nil, err
		}
		for i := 0; i < mag; i++ {
			row += dy
			col += dx
			steps++
			p.set(row, col, steps)
		}
	}
	return &p, nil
}

func (p *path) set(row, col, steps int) {
	rowEnt := p.rows[row]
	if rowEnt == nil {
		rowEnt = map[int]int{}
		p.rows[row] = rowEnt
	}
	if oldSteps, found := rowEnt[col]; found && oldSteps < steps {
		return
	}
	rowEnt[col] = steps
}

func (p path) isSet(row, col int) bool {
	if p.rows[row] == nil {
		return false
	}
	_, found := p.rows[row][col]
	return found
}

func (p path) steps(row, col int) int {
	if p.rows[row] == nil {
		return 0
	}
	return p.rows[row][col]
}

// Returns shortest path, path with fewest steps, error.
func (p *path) intersect(line string) (int, int, error) {
	distMin := 999999
	row := 0
	col := 0
	stepsMin := 999999
	steps := 0
	for _, entry := range strings.Split(line, ",") {
		dx, dy, mag, err := parseEntry(entry)
		if err != nil {
			return 0, 0, err
		}
		for i := 0; i < mag; i++ {
			row += dy
			col += dx
			steps += 1
			var dist int
			if p.isSet(row, col) {
				if row < 0 {
					dist = -row
				} else {
					dist = row
				}
				if col < 0 {
					dist -= col
				} else {
					dist += col
				}
				if steps + p.steps(row, col) < stepsMin {
					stepsMin = steps + p.steps(row, col)
				}
			}
			if dist == 0 {
				continue
			}
			if dist < distMin {
				distMin = dist
			}
		}			
	}
	return distMin, stepsMin, nil
}
