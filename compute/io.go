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

	NonBlocking bool

	Name string

	idle bool
	writeSem chan bool
}

func NewChanIO() (*ChanIO, *ChanIO) {
	a := make(chan int64, 30)
	b := make(chan int64, 30)
	// no buffering, in order to help catch errors
	return &ChanIO{
		Input:  a,
		Output: b,
		writeSem: make(chan bool, 1),
	}, &ChanIO{
		Input:  b,
		Output: a,
		writeSem: make(chan bool, 1),
	}
}

func (b *ChanIO) Write(val int64) error {
	b.writeSem <- true
	b.idle = false
	b.Output <- val
	<-b.writeSem
	return nil
}

func (b *ChanIO) WriteMulti(vals ...int64) error {
	b.writeSem <- true
	b.idle = false
	for _, val := range vals {
		b.Output <- val
	}
	<-b.writeSem
	return nil
}

func (b *ChanIO) Read() (int64, error) {
	t := 5*time.Second
	if b.NonBlocking {
		t = 1*time.Second
	}
	select {
	case res, ok := <-b.Input:
		b.idle = false
		if !ok {
			return 0, io.EOF
		}
		return res, nil
	case <-time.After(t):
		if b.NonBlocking {
			b.idle = true
			return -1, nil
		}
		return 0, fmt.Errorf("timeout in Read")
	}
}

func (b ChanIO) Close() {
	close(b.Output)
}

func (b ChanIO) Idle() bool {
	return b.idle
}
