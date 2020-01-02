package compute

import (
	"fmt"
	"io"
	"time"
)

type IO interface {
	Read() (int64, error)
	Write(val int64) error
	Close()
}

type BufIO struct {
	input  []int64
	inptr  int
	outbuf []int64
}

func (b *BufIO) Read() (int64, error) {
	if b.inptr >= len(b.input) {
		return 0, fmt.Errorf("read past end of input")
	}
	rc := b.input[b.inptr]
	b.inptr += 1
	return rc, nil
}

func (b *BufIO) Write(val int64) error {
	b.outbuf = append(b.outbuf, val)
	return nil
}

func (b BufIO) Output() []int64 {
	return b.outbuf
}

func (b BufIO) Close() {
}

func NewBufIO(inputs []int64) *BufIO {
	return &BufIO{
		input:  inputs,
		outbuf: []int64{},
	}
}

type ChanIO struct {
	Input  chan int64
	Output chan int64
}

func NewChanIO() (*ChanIO, *ChanIO) {
	a := make(chan int64)
	b := make(chan int64)
	// no buffering, in order to help catch errors
	return &ChanIO{
			Input:  a,
			Output: b,
		}, &ChanIO{
			Input:  b,
			Output: a,
		}
}

func (b ChanIO) Write(val int64) error {
	b.Output <- val
	return nil
}

func (b ChanIO) Read() (int64, error) {
	select {
	case res, ok := <-b.Input:
		if !ok {
			return 0, io.EOF
		}
		return res, nil
	case <-time.After(5 * time.Second):
		return 0, fmt.Errorf("timeout in Read")
	}
}

func (b ChanIO) Close() {
	close(b.Output)
}
