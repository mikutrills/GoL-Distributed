package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"uk.ac.bris.cs/gameoflife/stubs"
	//"uk.ac.bris.cs/gameoflife/util"
)

//process the GoL
// GOL WORKER

func main() {
	port := flag.String("port", "8030", "Port to listen on")
	flag.Parse()
	rpc.Register(&Ops{})
	listener, _ := net.Listen("tcp", ":"+*port)
	defer listener.Close()
	rpc.Accept(listener)
}

type Ops struct{}

func (g *Ops) WorldCalc(req stubs.Request, res *stubs.Response) (err error) {
	fmt.Println(req.Turns)
	turn := 0

	fmt.Println(req.World)

	res.World = req.World
	for turn < req.Turns {
		res.World = worldState(res.World, req)
		fmt.Println(res.World)
		turn++
	}
	fmt.Println("return: ")
	for row := range res.World {
		fmt.Println(res.World[row])
	}
	fmt.Println("return done ")

	//res.World = worldState(req.World, req)
	return
}

func worldState(finalWorld [][]uint8, req stubs.Request) [][]uint8 {
	tempWorld := world(req.ImageHeight, req.ImageWidth)
	wrapper := wrapperCalc(req.ImageHeight, req.ImageWidth)
	fmt.Println(wrapper)
	for i := 0; i < req.ImageHeight; i++ {
		for j := 0; j < req.ImageWidth; j++ {
			AliveCellsCount := aliveCount(i, j, finalWorld, wrapper)
			worldStateCalc(finalWorld, tempWorld, i, j, AliveCellsCount)

		}
	}

	return tempWorld

}

//worl size 16, 0xF
//world size 64 0x3F
//orld sze 512 0x1FF
// Optimization for +16 binary instead as -1 = 15 in binary
func aliveCount(i int, j int, tempWorld [][]uint8, wrapper int) uint8 {
	//check live neighbors
	var count uint8
	for y := j - 1; y <= j+1; y++ {
		for x := i - 1; x <= i+1; x++ {
			if !(x == i && y == j) {
				//tempWorld[(x+wrapper)%len(tempWorld)][(y+wrapper)%len(tempWorld[0])]
				if tempWorld[x&wrapper][y&wrapper] == 255 {
					count++
				}
			}
		}

	}
	return count
}

func world(ImageHeight int, ImageWidth int) [][]uint8 {
	worldSlice := make([][]uint8, ImageHeight)
	for i := range worldSlice {
		worldSlice[i] = make([]uint8, ImageWidth)
	}
	return worldSlice
}

func wrapperCalc(ImageHeight int, ImageWidth int) int {
	var wrapper int
	if ImageWidth == 16 && ImageHeight == 16 {
		wrapper = 0xF

	} else if ImageWidth == 64 && ImageHeight == 64 {
		wrapper = 0x3F
	} else if ImageWidth == 512 && ImageHeight == 512 {
		wrapper = 0x1FF

	}
	return wrapper
}

func worldStateCalc(finalWorld [][]uint8, tempWorld [][]uint8, i int, j int, AliveCellsCount uint8) {
	if finalWorld[i][j] == 255 {

		if AliveCellsCount > 3 {
			tempWorld[i][j] = 0
		} else if AliveCellsCount < 2 {
			tempWorld[i][j] = 0
		} else if AliveCellsCount == 2 || AliveCellsCount == 3 {
			tempWorld[i][j] = 255
		} else {
			tempWorld[i][j] = 255
		}

	}

	if finalWorld[i][j] == 0 {
		if AliveCellsCount == 3 {
			tempWorld[i][j] = 255
		} else {
			tempWorld[i][j] = 0
		}
	}

}
