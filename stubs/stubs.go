package stubs

var GolCalc = "Ops.WorldCalc"

// Response Returns board state to local controller
type Response struct {
	World [][]uint8
}

// Request
type Request struct {
	World       [][]uint8
	ImageHeight int
	ImageWidth  int
	Turns       int
}
