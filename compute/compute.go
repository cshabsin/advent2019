package compute

import "fmt"

func sizeUp(buffer []int, index int) []int {
	for index >= len(buffer) {
		buffer = append(buffer, 0)
	}
	return buffer
}

type Intcode struct {
	memory []int
	ip     int

	input []int
	inptr int

	outbuf []int
}

func NewIntcode(inbuf, inputs []int) *Intcode {
	// copy the buffer to leave original unchanged.
	buffer := make([]int, 0, len(inbuf))
	for _, b := range inbuf {
		buffer = append(buffer, b)
	}
	return &Intcode{
		memory: buffer,
		input:  inputs,
		outbuf: []int{},
	}
}

func (i Intcode) opcode() int {
	return i.memory[i.ip] % 100
}

func (i Intcode) modes() modes {
	return makeModes(i.memory[i.ip])
}

func (i Intcode) evalParam(index int) (int, error) {
	return makeModes(i.memory[i.ip]).evalParam(index, i.memory, i.memory[i.ip+1+index])
}

func (i *Intcode) setMemoryByParam(index, value int) {
	target := i.memory[i.ip+1+index]
	fmt.Printf("target: %d\n", target)
	i.memory = sizeUp(i.memory, target)
	i.memory[target] = value
}

func (i *Intcode) read() (int, error) {
	if i.inptr >= len(i.input) {
		return 0, fmt.Errorf("read past end of input")
	}
	rc := i.input[i.inptr]
	i.inptr += 1
	return rc, nil
}

func (i *Intcode) output(val int) {
	i.outbuf = append(i.outbuf, val)
}

func (i *Intcode) jump(newIp int) {
	i.ip = newIp
}

func (i *Intcode) Run() ([]int, []int, error) {
	for {
		var opwidth int
		switch i.opcode() {
		case 1:
			// ADD
			a, err := i.evalParam(0)
			if err != nil {
				return nil, nil, err
			}
			b, err := i.evalParam(1)
			if err != nil {
				return nil, nil, err
			}
			i.setMemoryByParam(2, a+b)
			opwidth = 4
		case 2:
			// MUL
			a, err := i.evalParam(0)
			if err != nil {
				return nil, nil, err
			}
			b, err := i.evalParam(1)
			if err != nil {
				return nil, nil, err
			}
			i.setMemoryByParam(2, a*b)
			opwidth = 4
		case 3:
			// IN
			val, err := i.read()
			if err != nil {
				return nil, nil, fmt.Errorf("IN: %v at IP %d", err, i.ip)
			}
			i.setMemoryByParam(0, val)
			opwidth = 2
		case 4:
			// OUT
			a, err := i.evalParam(0)
			if err != nil {
				return nil, nil, err
			}
			i.output(a)
			opwidth = 2
		case 5:
			// JT (jump if true)
			opwidth = 3
			a, err := i.evalParam(0)
			if err != nil {
				return nil, nil, err
			}
			if a != 0 {
				b, err := i.evalParam(1)
				if err != nil {
					return nil, nil, err
				}
				i.jump(b)
				opwidth = 0
			}
		case 6:
			// JF (jump if false)
			opwidth = 3
			a, err := i.evalParam(0)
			if err != nil {
				return nil, nil, err
			}
			if a == 0 {
				b, err := i.evalParam(1)
				if err != nil {
					return nil, nil, err
				}
				i.jump(b)
				opwidth = 0
			}
		case 7:
			// LT (less than)
			a, err := i.evalParam(0)
			if err != nil {
				return nil, nil, err
			}
			b, err := i.evalParam(1)
			if err != nil {
				return nil, nil, err
			}
			if a < b {
				i.setMemoryByParam(2, 1)
			} else {
				i.setMemoryByParam(2, 0)
			}
			opwidth = 4
		case 8:
			// EQ (equal)
			a, err := i.evalParam(0)
			if err != nil {
				return nil, nil, err
			}
			b, err := i.evalParam(1)
			if err != nil {
				return nil, nil, err
			}
			if a == b {
				i.setMemoryByParam(2, 1)
			} else {
				i.setMemoryByParam(2, 0)
			}
			opwidth = 4
		case 99:
			// EXIT
			return i.memory, i.outbuf, nil
		default:
			return nil, nil, fmt.Errorf("Invalid opcode %d at IP %d", i.memory[i.ip], i.ip)
		}
		i.ip += opwidth
	}
}
