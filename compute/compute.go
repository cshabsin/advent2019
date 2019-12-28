package compute

import "fmt"

func sizeUp(buffer []int, index int) []int {
	for index >= len(buffer) {
		buffer = append(buffer, 0)
	}
	return buffer
}

func Run(inbuf, inputs []int) ([]int, []int, error) {
	// copy the buffer to leave original unchanged.
	buffer := make([]int, 0, len(inbuf))
	for _, b := range inbuf {
		buffer = append(buffer, b)
	}
	outbuf := make([]int, 0)
	ip := 0
	inptr := 0
	for {
		var opwidth int
		opcode := buffer[ip] % 100
		modes := makeModes(buffer[ip])
		switch opcode {
		case 1:
			// ADD
			a, err := modes.evalParam(0, buffer, buffer[ip+1])
			if err != nil {
				return nil, nil, err
			}
			b, err := modes.evalParam(1, buffer, buffer[ip+2])
			if err != nil {
				return nil, nil, err
			}
			buffer = sizeUp(buffer, buffer[ip+3])
			buffer[buffer[ip+3]] = a + b
			opwidth = 4
		case 2:
			// MUL
			a, err := modes.evalParam(0, buffer, buffer[ip+1])
			if err != nil {
				return nil, nil, err
			}
			b, err := modes.evalParam(1, buffer, buffer[ip+2])
			if err != nil {
				return nil, nil, err
			}
			buffer = sizeUp(buffer, buffer[ip+3])
			buffer[buffer[ip+3]] = a * b
			opwidth = 4
		case 3:
			// SAV
			if inptr >= len(inputs) {
				return nil, nil, fmt.Errorf("SAV: Read past end of input at IP %d", ip)
			}
			val := inputs[inptr]
			inptr++
			buffer = sizeUp(buffer, buffer[ip+1])
			buffer[buffer[ip+1]] = val
			opwidth = 2
		case 4:
			// OUT
			a, err := modes.evalParam(0, buffer, buffer[ip+1])
			if err != nil {
				return nil, nil, err
			}
			outbuf = append(outbuf, a)
			opwidth = 2
		case 5:
			// JT (jump if true)
			opwidth = 3
			a, err := modes.evalParam(0, buffer, buffer[ip+1])
			if err != nil {
				return nil, nil, err
			}
			if a != 0 {
				b, err := modes.evalParam(1, buffer, buffer[ip+2])
				if err != nil {
					return nil, nil, err
				}
				ip = b
				opwidth = 0
			}
		case 6:
			// JF (jump if false)
			opwidth = 3
			a, err := modes.evalParam(0, buffer, buffer[ip+1])
			if err != nil {
				return nil, nil, err
			}
			if a == 0 {
				b, err := modes.evalParam(1, buffer, buffer[ip+2])
				if err != nil {
					return nil, nil, err
				}
				ip = b
				opwidth = 0
			}
		case 7:
			// LT (less than)
			a, err := modes.evalParam(0, buffer, buffer[ip+1])
			if err != nil {
				return nil, nil, err
			}
			b, err := modes.evalParam(1, buffer, buffer[ip+2])
			if err != nil {
				return nil, nil, err
			}
			buffer = sizeUp(buffer, buffer[ip+3])
			if a < b {
				buffer[buffer[ip+3]] = 1
			} else {
				buffer[buffer[ip+3]] = 0
			}
			opwidth = 4
		case 8:
			// EQ (equal)
			a, err := modes.evalParam(0, buffer, buffer[ip+1])
			if err != nil {
				return nil, nil, err
			}
			b, err := modes.evalParam(1, buffer, buffer[ip+2])
			if err != nil {
				return nil, nil, err
			}
			buffer = sizeUp(buffer, buffer[ip+3])
			if a == b {
				buffer[buffer[ip+3]] = 1
			} else {
				buffer[buffer[ip+3]] = 0
			}
			opwidth = 4
		case 99:
			// EXIT
			return buffer, outbuf, nil
		default:
			return nil, nil, fmt.Errorf("Invalid opcode %d at IP %d", buffer[ip], ip)
		}
		ip += opwidth
	}
}
