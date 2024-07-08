package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/davlaf/minimax-connect-4/src/board"
)

func main() {
	// game_board := board.Board{
	// 	XPlayer: 0b100101001110001011000010000100000010001101,
	// 	OPlayer: 0b11010010001000100001101000011000101000010,
	// 	TurnNumber: 0,
	// }
	game_board := board.Board{
		XPlayer: 0b0,
		OPlayer: 0b0,
		TurnNumber: 0,
	}
	// game_board.Print()
	// fmt.Println("----------- enter number 1 to 7")
	
	list_of_actions_per_cpu, number_of_tasks := game_board.GetActionsForNumberOfCPU()
	evaluated_actions_channel := make(chan board.EvaluatedAction, number_of_tasks)
	var searched_games uint64 = 0
	var wg sync.WaitGroup
	for cpu_index, cpu_action_list := range list_of_actions_per_cpu {
		fmt.Printf("cpu_index: %v\n", cpu_index)
		board.PrintBoardList(cpu_action_list)
		wg.Add(1)
		go func() {
			for _, action := range cpu_action_list {
				defer wg.Done()
				evaluated_actions_channel <- action.FindBestNextMove(&searched_games)
			}
		}()
	}
	fmt.Println(time.Now())
	wg.Wait()
	var evaluated_actions []board.EvaluatedAction
	for evaluated_action := range evaluated_actions_channel {
		evaluated_actions = append(evaluated_actions, evaluated_action)
	}

	board.PrintActionList(evaluated_actions)
	
	
	
	// fmt.Printf("searched_games: %v\n", searched_games)
	// best_move.Print()
	
	// for {
	// 	user_move := game_board.GetNextMoveFromInput()
	// 	game_board = user_move
	// 	game_board.Print()
	// 	fmt.Println("----------- enter number 1 to 7")
	// 	is_terminal, value := game_board.IsTerminalWithValue()
	// 	if is_terminal {
	// 		switch value {
	// 			case 1: fmt.Println("# wins!")
	// 			case 0: fmt.Println("Draw!")
	// 			case -1: fmt.Println("/ wins!")
	// 		}
	// 		break
	// 	}
	// 	// fmt.Printf("game_board.XPlayer: %b\n", game_board.XPlayer)
	// 	// fmt.Printf("game_board.OPlayer: %b\n", game_board.OPlayer)
	// }
	
}