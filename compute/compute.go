package compute

import (
	"fmt"
	"strconv"
	"strings"
)

func sizeUp(buffer []int64, index int64) []int64 {
	for index >= int64(len(buffer)) {
		buffer = append(buffer, 0)
	}
	return buffer
}

type Intcode struct {
	memory  []int64
	ip      int
	io      IO
	relbase int64
}

func ParseFile(content []byte) ([]int64, error) {
	vals := strings.Split(string(content), ",")
	buf := make([]int64, 0, len(vals))
	for _, v := range vals {
		intval, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return nil, err
		}
		buf = append(buf, intval)
	}
	return buf, nil
}

func NewIntcode(inbuf []int64, io IO) *Intcode {
	// copy the buffer to leave original unchanged.
	buffer := make([]int64, 0, len(inbuf))
	for _, b := range inbuf {
		buffer = append(buffer, b)
	}
	return &Intcode{
		memory: buffer,
		io:     io,
	}
}

func (i Intcode) opcode() int64 {
	return i.memory[i.ip] % 100
}

func (i Intcode) modes() modes {
	return makeModes(i.memory[i.ip])
}

func (i Intcode) evalParam(index int) (int64, error) {
	return makeModes(i.memory[i.ip]).evalParam(index, i.memory, i.memory[i.ip+1+index], i.relbase)
}

func (i *Intcode) setMemoryByParam(index int, value int64) {
	m := makeModes(i.memory[i.ip])
	var target int64
	switch m.getMode(index) {
	case 0:
		target = i.memory[i.ip+1+index]
	case 2:
		target = i.relbase + i.memory[i.ip+1+index]
	}
	i.memory = sizeUp(i.memory, target)
	i.memory[target] = value
}

func (i *Intcode) read() (int64, error) {
	return i.io.Read()
}

func (i *Intcode) output(val int64) error {
	return i.io.Write(val)
}

func (i *Intcode) jump(newIp int) {
	i.ip = newIp
}

func (i *Intcode) Run() ([]int64, error) {
	for {
		var opwidth int
		switch i.opcode() {
		case 1:
			// ADD
			a, err := i.evalParam(0)
			if err != nil {
				return nil, err
			}
			b, err := i.evalParam(1)
			if err != nil {
				return nil, err
			}
			i.setMemoryByParam(2, a+b)
			opwidth = 4
		case 2:
			// MUL
			a, err := i.evalParam(0)
			if err != nil {
				return nil, err
			}
			b, err := i.evalParam(1)
			if err != nil {
				return nil, err
			}
			i.setMemoryByParam(2, a*b)
			opwidth = 4
		case 3:
			// IN
			val, err := i.read()
			if err != nil {
				return nil, fmt.Errorf("IN: %w at IP %d", err, i.ip)
			}
			i.setMemoryByParam(0, val)
			opwidth = 2
		case 4:
			// OUT
			a, err := i.evalParam(0)
			if err != nil {
				return nil, err
			}
			i.output(a)
			opwidth = 2
		case 5:
			// JT (jump if true)
			opwidth = 3
			a, err := i.evalParam(0)
			if err != nil {
				return nil, err
			}
			if a != 0 {
				b, err := i.evalParam(1)
				if err != nil {
					return nil, err
				}
				i.jump(int(b))
				opwidth = 0
			}
		case 6:
			// JF (jump if false)
			opwidth = 3
			a, err := i.evalParam(0)
			if err != nil {
				return nil, err
			}
			if a == 0 {
				b, err := i.evalParam(1)
				if err != nil {
					return nil, err
				}
				i.jump(int(b))
				opwidth = 0
			}
		case 7:
			// LT (less than)
			a, err := i.evalParam(0)
			if err != nil {
				return nil, err
			}
			b, err := i.evalParam(1)
			if err != nil {
				return nil, err
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
				return nil, err
			}
			b, err := i.evalParam(1)
			if err != nil {
				return nil, err
			}
			if a == b {
				i.setMemoryByParam(2, 1)
			} else {
				i.setMemoryByParam(2, 0)
			}
			opwidth = 4
		case 9:
			// RELBASE
			a, err := i.evalParam(0)
			if err != nil {
				return nil, err
			}
			i.relbase += a
			opwidth = 2
		case 99:
			// EXIT
			i.io.Close()
			return i.memory, nil
		default:
			return nil, fmt.Errorf("Invalid opcode %d at IP %d", i.memory[i.ip], i.ip)
		}
		i.ip += opwidth
	}
}
