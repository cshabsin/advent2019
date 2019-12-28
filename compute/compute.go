package compute

import "fmt"

func Run(instructions []int) ([]int, error) {
	buffer := make([]int, 0, len(instructions))
	for _, b := range instructions {
		buffer = append(buffer, b)
	}
	ip := 0
	for {
		var opwidth int
		opcode := buffer[ip] % 100
		modes := makeModes(buffer[ip])
		switch opcode {
		case 1:
			// ADD
			a, err := modes.evalParam(0, buffer, buffer[ip+1])
			if err != nil {
				return nil, err
			}
			b, err := modes.evalParam(1, buffer, buffer[ip+2])
			if err != nil {
				return nil, err
			}
			buffer[buffer[ip+3]] = a + b
			opwidth = 4
		case 2:
			// MUL
			a, err := modes.evalParam(0, buffer, buffer[ip+1])
			if err != nil {
				return nil, err
			}
			b, err := modes.evalParam(1, buffer, buffer[ip+2])
			if err != nil {
				return nil, err
			}
			buffer[buffer[ip+3]] = a * b
			opwidth = 4
		case 99:
			// EXIT
			return buffer, nil
		default:
			return nil, fmt.Errorf("Invalid opcode %d at IP %d", buffer[ip], ip)
		}
		ip += opwidth
	}
}
