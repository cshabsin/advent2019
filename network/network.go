package network

import (
	"fmt"
	"time"

	"github.com/cshabsin/advent2019/compute"
)

type Network struct {
	computers []computer
	nIOs      []*compute.ChanIO
	stage     []int

	nat *nat
}

func NewNetwork() *Network {
	n := &Network{}
	n.nat = &nat{
		activityChan: make(chan bool),
		n:            n,
	}
	return n
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
	go n.nat.monitor()
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
		n.nat.tickle()
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

		//		fmt.Printf("(%d)sending to %d: %d, %d\n", id, addr, x, y)

		if addr == 255 {
			n.nat.update(x, y)
			continue
		}

		if int(addr) >= len(n.computers) {
			fmt.Printf("invalid address %d\n", addr)
			continue
		}

		if err := n.SendTo(int(addr), x, y); err != nil {
			fmt.Printf("WriteMulti(%d) x: %v\n", err)
			return
		}
	}
}

func (n *Network) SendTo(addr int, x, y int64) error {
	return n.nIOs[addr].WriteMulti(x, y)
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

func (c computer) isIdle() bool {
	return c.io.Idle()
}

type nat struct {
	x, y         int64
	prevX, prevY int64

	activityChan chan bool
	n            *Network
}

// update returns true if the values are unchanged
func (nat *nat) update(x, y int64)  {
	nat.x, nat.y = x, y
}

func (nat *nat) tickle() {
	nat.activityChan <- true
}

func (nat *nat) monitor() {
	<-nat.activityChan
	for {
		select {
		case <-nat.activityChan:
		case <-time.After(100 * time.Millisecond):
			idle := true
			for _, c := range nat.n.computers {
				if !c.isIdle() {
					idle = false
					break
				}
			}
			if idle {
				fmt.Printf("NAT activated: %d, %d\n", nat.x, nat.y)
				nat.n.SendTo(0, nat.x, nat.y)
				if nat.x == nat.prevX && nat.y == nat.prevY {
					fmt.Println("Repeated values!")
				}
				nat.prevX, nat.prevY = nat.x, nat.y
			}
		}
	}
}
