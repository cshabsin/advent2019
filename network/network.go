package network

import (
	"fmt"

	"github.com/cshabsin/advent2019/compute"
)

type Network struct {
	computers []computer
	nIOs      []*compute.ChanIO
	stage []int
}

func NewNetwork() *Network {
	return &Network{}
}

func (n *Network) AddComputer(buf []int64, id int) {
	cIO, nIO := compute.NewChanIO()
	cIO.NonBlocking = true
	cIO.Name = fmt.Sprintf("c%d", id)
	n.computers = append(n.computers, computer{
		id:  id,
		buf: buf,
		io:  cIO,
	})

	nIO.Name = fmt.Sprintf("n%d", id)
	n.nIOs = append(n.nIOs, nIO)

	n.stage = append(n.stage, 0)
}

func (n *Network) Run() {
	fin := make(chan bool)
	for i, c := range n.computers {
		i, c := i, c
		go c.run(fin)
		go n.dispatch(i)
	}
	for _ = range n.computers {
		<-fin
	}
}

func (n *Network) dispatch(id int) {
	nIO := n.nIOs[id]
	if err := nIO.Write(int64(id)); err != nil {
		fmt.Printf("Write(assign addr %d): %v\n", id, err)
		return
	}
	for {
		nIO.NonBlocking = true
		var addr int64
		for {
			var err error
			addr, err = nIO.Read()
			if err != nil {
				fmt.Printf("Read(%d) addr: %v\n", id, err)
				return
			}
			if addr != -1 {
				break
			}
		}
		nIO.NonBlocking = false
		x, err := nIO.Read()
		if err != nil {
			fmt.Printf("Read(%d) x (->%d): %v\n", id, addr, err)
			continue
		}
		y, err := nIO.Read()
		if err != nil {
			fmt.Printf("Read(%d) y (->%d): %v\n", id, addr, err)
			continue
		}

		fmt.Printf("(%d)sending to %d: %d, %d\n", id, addr, x, y)
		
		if addr == 255 {
			fmt.Printf("sent to 255: %d, %d\n", x, y)
		}
		
		if int(addr) >= len(n.computers) {
			continue
		}

		nIO := n.nIOs[int(addr)]
		
		if err := nIO.WriteMulti(x, y); err != nil {
			fmt.Printf("WriteMulti(%d) x: %v\n", err)
			return
		}
	}
}

type computer struct {
	id  int
	buf []int64
	io  *compute.ChanIO
}

func (c *computer) run(fin chan bool) {
	intcode := compute.NewIntcode(c.buf, c.io)
	if _, err := intcode.Run(); err != nil {
		fmt.Printf("Run(%d): %v\n", c.id, err)
	}
	fin <- true
}
