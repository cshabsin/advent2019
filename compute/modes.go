package compute

import "fmt"

// Info about the parameter modes of the parameters of an opcode
type modes []int

func makeModes(opcode int64) modes {
	opcode /= 100
	rc := modes{}
	for {
		if opcode == 0 {
			return rc
		}
		rc = append(rc, int(opcode%10))
		opcode /= 10
	}
}

// evaluates param based on mode
// i - param index
// b - memory buffer of computer
// p - param value to evaluate (position or immediate)
func (m modes) evalParam(i int, b []int64, p, relbase int64) (int64, error) {
	mode := m.getMode(i)
	switch mode {
	case 0:
		if p >= int64(len(b)) {
			return 0, nil
		}
		return b[p], nil
	case 1:
		return p, nil
	case 2:
		if p+relbase >= int64(len(b)) {
			return 0, nil
		}
		return b[p+relbase], nil
	default:
		return 0, fmt.Errorf("evalParam: invalid mode %d", mode)
	}
}

func (m modes) getMode(index int) int {
	if index < len(m) {
		return m[index]
	}
	return 0
}
