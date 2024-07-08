package main

import "fmt"

type BoardBits struct {
	Bits uint64
}

func (bits *BoardBits) MoveLeft(amount int) {
	bits.Bits = bits.Bits << (6*uint64(amount))
}

func (bits *BoardBits) MoveUp(amount int) {
	bits.Bits = bits.Bits << (1*uint64(amount))
}

func (bits BoardBits) Print() {
	var lines [6]uint64
	for i := 6; i >= 0; i-- {
		column := bits.Bits >> uint64(6*i)
		lines[0] = lines[0] << 1 | column & 0b100000 >> 5
		lines[1] = lines[1] << 1 | column & 0b010000 >> 4
		lines[2] = lines[2] << 1 | column & 0b001000 >> 3
		lines[3] = lines[3] << 1 | column & 0b000100 >> 2
		lines[4] = lines[4] << 1 | column & 0b000010 >> 1
		lines[5] = lines[5] << 1 | column & 0b000001 >> 0
	}

	for i := 0; i < 6; i++ {
		fmt.Printf("%07b\n", lines[i])
	}
}

func main2() {
	horizontal_board := BoardBits{
		Bits: 0b000001000001000001000001,
	}
	vertical_board := BoardBits{
		Bits: 0b1111,
	}
	diagonal_left := BoardBits{
		Bits: 0b001000000100000010000001,
	}
	diagonal_right := BoardBits{
		Bits: 0b000001000010000100001000,
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 6; j++ {
			board := horizontal_board
			board.MoveLeft(i)
			board.MoveUp(j)
			fmt.Printf("0b%b\n",board.Bits)
		}
	}
	for i := 0; i < 7; i++ {
		for j := 0; j < 3; j++ {
			board := vertical_board
			board.MoveLeft(i)
			board.MoveUp(j)
			fmt.Printf("0b%b\n",board.Bits)
		}
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			board := diagonal_left
			board.MoveLeft(i)
			board.MoveUp(j)
			fmt.Printf("0b%b\n",board.Bits)
		}
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			board := diagonal_right
			board.MoveLeft(i)
			board.MoveUp(j)
			fmt.Printf("0b%b\n",board.Bits)
		}
	}
}

func GenerateWinningMoves() {

}