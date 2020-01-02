package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/cshabsin/advent2019/compute"
)

func main() {
	content, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	buf, err := compute.ParseFile(content)
	if err != nil {
		log.Fatal(err)
	}
	mapper := droneCtl{buf}
	mapper.partOne()
}

func (d droneCtl) partOne() {
	count := 0
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			moving, err := d.sendDrone(x, y)
			if err != nil {
				fmt.Printf("at (%d, %d): %v\n", x, y, err)
				return
			}
			if moving {
				fmt.Printf("#")
				count++
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("count: %d\n", count)
}

type droneCtl struct {
	buf []int64 // program input
}

func (d droneCtl) sendDrone(x, y int) (bool, error) {
	io := compute.NewBufIO([]int64{int64(x), int64(y)})
	intcode := compute.NewIntcode(d.buf, io)
	if _, err := intcode.Run(); err != nil {
		return false, err
	}
	if len(io.Output()) != 1 {
		return false, fmt.Errorf("unexpected output buffer: %v", io.Output())
	}
	return io.Output()[0] == 1, nil
}

func RunGrid(aIO compute.IO) error {
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			if err := aIO.Write(int64(x)); err != nil {
				return err
			}
			if err := aIO.Write(int64(y)); err != nil {
				return err
			}
			// val, err := aIO.Read()
			// if err != nil {
			// 	return err
			// }
			// fmt.Printf("%d, %d: %d\n", x, y, val)
		}
	}
	return nil
}
