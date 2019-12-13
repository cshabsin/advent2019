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
		switch buffer[ip] {
		case 1:
			// ADD
			a := buffer[buffer[ip+1]]
			b := buffer[buffer[ip+2]]
			buffer[buffer[ip+3]] = a + b
			opwidth = 4
		case 2:
			// MUL
			a := buffer[buffer[ip+1]]
			b := buffer[buffer[ip+2]]
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
