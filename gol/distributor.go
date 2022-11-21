package gol

import (
	"fmt"
	"net/rpc"
	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
}

//client
//GoL engine
//transfer GoL to server
// distributor divides the work between workers and interacts with other goroutines.

func distributor(p Params, c distributorChannels) {
	// TODO: Create a 2D slice to store the world.
	c.ioCommand <- ioInput
	fp := fmt.Sprintf("%vx%v", p.ImageWidth, p.ImageHeight)
	c.ioFilename <- fp
	finalWorld := world(p)

	turn := 0

	// Populating world with image
	for x := 0; x < p.ImageWidth; x++ {
		for y := 0; y < p.ImageHeight; y++ {
			finalWorld[x][y] = <-c.ioInput
		}
	}

	// TODO: Execute all turns of the Game of Life.
	//Client code
	var ip = "127.0.0.1:8030"
	fmt.Println("Server: ", ip)
	client, err := rpc.Dial("tcp", ip)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer client.Close()

	request := stubs.Request{World: finalWorld, ImageHeight: p.ImageHeight, ImageWidth: p.ImageWidth, Turns: p.Turns}
	response := new(stubs.Response)
	client.Call(stubs.GolCalc, request, response)

	// TODO: Report the final state using FinalTurnCompleteEvent.
	aliveSlice := calcAlive(p, response.World)
	c.events <- FinalTurnComplete{turn, aliveSlice}

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle
	c.events <- StateChange{turn, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)

}

func world(p Params) [][]uint8 {
	worldSlice := make([][]uint8, p.ImageHeight)
	for i := range worldSlice {
		worldSlice[i] = make([]uint8, p.ImageWidth)
	}
	return worldSlice
}

func cellFlipped(c distributorChannels, x int, y int) {
	c.events <- CellFlipped{CompletedTurns: 0, Cell: util.Cell{x, y}}

}

func calcAlive(p Params, world [][]uint8) []util.Cell {
	var cells []util.Cell
	for i := 0; i < p.ImageHeight; i++ {
		for j := 0; j < p.ImageWidth; j++ {
			if world[i][j] == 255 {
				cells = append(cells, util.Cell{X: j, Y: i})
			}
		}
	}

	return cells

}
