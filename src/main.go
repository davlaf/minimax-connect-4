package main

import (
	"fmt"

	"github.com/davlaf/minimax-connect-4/src/board"
)

func main() {
	// game_board := board.Board{
	// 	XPlayer: 0b11000001001000000001000110000110000001,
	// 	OPlayer: 0b10000111001110000001000001000000,
	// 	TurnNumber: 19,
	// }
	game_board := board.Board{
		XPlayer: 0b0,
		OPlayer: 0b0,
		TurnNumber: 0,
	}
	game_board.Print()
	fmt.Println("----------- enter number 1 to 7")
	
	
	
	// fmt.Println(time.Now())
	
	// fmt.Printf("searched_games: %v\n", searched_games)
	// best_move.Print()
	
	// game_board.Print()
	// board_value, depth := game_board.Minimax(42, -128, 127, &searched_games)
	// board.EvaluatedAction{
	// 	Action: game_board,
	// 	ActionValue: board_value,
	// 	TerminationDepth: depth,
	// }.Print()
	// game_board.FindBestNextMove(&searched_games)


	var is_terminal bool
	var value int8
	for {
		bot_move, considered_moves := game_board.FindBestNextMove(12)
		board.PrintActionList(considered_moves)
		game_board = bot_move
		game_board.Print()
		is_terminal, value = game_board.IsTerminalWithValue()
		if is_terminal {
			switch value {
				case 9: fmt.Println("# wins!")
				case 0: fmt.Println("Draw!")
				case -9: fmt.Println("/ wins!")
			}
			break
		}
		user_move := game_board.GetNextMoveFromInput()
		game_board = user_move
		game_board.Print()
		fmt.Println("----------- enter number 1 to 7")
		is_terminal, value = game_board.IsTerminalWithValue()
		if is_terminal {
			switch value {
				case 9: fmt.Println("# wins!")
				case 0: fmt.Println("Draw!")
				case -9: fmt.Println("/ wins!")
			}
			break
		}
		fmt.Printf("game_board.XPlayer: %b\n", game_board.XPlayer)
		fmt.Printf("game_board.OPlayer: %b\n", game_board.OPlayer)
		fmt.Printf("game_board.TurnNumber: %v\n", game_board.TurnNumber)
		
	}
	
}