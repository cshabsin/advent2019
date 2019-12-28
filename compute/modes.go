package compute

import "fmt"

// Info about the parameter modes of the parameters of an opcode
type modes []int

func makeModes(opcode int) modes {
	opcode /= 100
	rc := modes{}
	for {
		if opcode == 0 {
			return rc
		}
		rc = append(rc, opcode % 10)
		opcode /= 10
	}
}

// evaluates param based on mode
// i - param index
// b - memory buffer of computer
// p - param value to evaluate (position or immediate)
func (m modes) evalParam(i int, b []int, p int) (int, error) {
	var mode int
	if i < len(m) {
		mode = m[i]
	}
	if mode == 0 {
		return b[p], nil
	}
	if mode == 1 {
		return p, nil
	}
	return 0, fmt.Errorf("evalParam: invalid mode %d", mode)
}
