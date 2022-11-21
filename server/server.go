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

	res.World = req.World
	for turn < req.Turns {
		res.World = worldState(res.World, req)
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
	fmt.Println("hello")
	tempWorld := world(req.ImageHeight, req.ImageWidth)
	//tempWorld = finalWorld
	//goes through rows
	wrapper := wrapperCalc(req.ImageHeight, req.ImageWidth)
	for i, r := range finalWorld {
		//goes through each cell in row
		for j, color := range r {
			AliveCellsCount := aliveCount(i, j, finalWorld, wrapper)
			var newColor uint8

			if color == 255 {

				if AliveCellsCount > 3 {
					newColor = 0
				} else if AliveCellsCount < 2 {
					newColor = 0
				} else if AliveCellsCount == 2 || AliveCellsCount == 3 {
					newColor = 255
				} else {
					newColor = 255
				}

			}

			if color == 0 {
				if AliveCellsCount == 3 {
					newColor = 255
				} else {
					newColor = 0
				}
			}

			//if color != newColor {
			//fmt.Println("Shit the bed")
			//}

			tempWorld[i][j] = newColor

		}
	}
	finalWorld = tempWorld

	for i := 0; i < req.ImageHeight; i++ {
		for j := 0; j < req.ImageWidth; j++ {
			if tempWorld[i][j] != finalWorld[i][j] {
				fmt.Println(i, j, "help")
			}
		}
	}

	return finalWorld

}

//worl size 16, 0xF
//world size 64 0x3F
//orld sze 512 0x1FF
// Optimization for +16 binary instead as -1 = 15 in binary
func aliveCount(i int, j int, tempWorld [][]uint8, wrapper int) uint8 {
	//check live neighbors
	var count uint8
	for x := i - 1; x <= i+1; x++ {
		for y := j - 1; y <= j+1; y++ {
			if !(x == i && y == j) {
				//if tempWorld[(x+16)%len(tempWorld)][(y+16)%len(tempWorld)] == 255 {
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
