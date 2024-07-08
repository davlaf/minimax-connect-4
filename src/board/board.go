package board

import (
	"fmt"
	"math"
	"math/bits"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

const (          
	FULL_BOARD   = 0b111111111111111111111111111111111111111111
	FULL_COLUMN  = 0b111111
	O_MIN_PLAYER = 1
	X_MAX_PLAYER = 0
)

type Board struct {
	XPlayer    uint64 // max player
	OPlayer    uint64 // min player
	TurnNumber uint8  // first turn is 0 and X starts
}

// 42 36 30 24 18 12  6
// 41 35 29 23 17 11  5
// 40 34 28 22 16 10  4 
// 39 33 27 21 15  9  3
// 38 32 26 20 14  8  2
// 37 31 25 19 13  7  1
// MSB                          LSB
// - - ... - - - 42 41 40 ... 3 2 1

func (board Board) IsTerminalWithValue() (bool, int8) {
	winning_moves := [...]uint64{
		0b1000001000001000001,
		0b10000010000010000010,
		0b100000100000100000100,
		0b1000001000001000001000,
		0b10000010000010000010000,
		0b100000100000100000100000,
		0b1000001000001000001000000,
		0b10000010000010000010000000,
		0b100000100000100000100000000,
		0b1000001000001000001000000000,
		0b10000010000010000010000000000,
		0b100000100000100000100000000000,
		0b1000001000001000001000000000000,
		0b10000010000010000010000000000000,
		0b100000100000100000100000000000000,
		0b1000001000001000001000000000000000,
		0b10000010000010000010000000000000000,
		0b100000100000100000100000000000000000,
		0b1000001000001000001000000000000000000,
		0b10000010000010000010000000000000000000,
		0b100000100000100000100000000000000000000,
		0b1000001000001000001000000000000000000000,
		0b10000010000010000010000000000000000000000,
		0b100000100000100000100000000000000000000000,
		0b1111,
		0b11110,
		0b111100,
		0b1111000000,
		0b11110000000,
		0b111100000000,
		0b1111000000000000,
		0b11110000000000000,
		0b111100000000000000,
		0b1111000000000000000000,
		0b11110000000000000000000,
		0b111100000000000000000000,
		0b1111000000000000000000000000,
		0b11110000000000000000000000000,
		0b111100000000000000000000000000,
		0b1111000000000000000000000000000000,
		0b11110000000000000000000000000000000,
		0b111100000000000000000000000000000000,
		0b1111000000000000000000000000000000000000,
		0b11110000000000000000000000000000000000000,
		0b111100000000000000000000000000000000000000,
		0b1000000100000010000001,
		0b10000001000000100000010,
		0b100000010000001000000100,
		0b1000000100000010000001000000,
		0b10000001000000100000010000000,
		0b100000010000001000000100000000,
		0b1000000100000010000001000000000000,
		0b10000001000000100000010000000000000,
		0b100000010000001000000100000000000000,
		0b1000000100000010000001000000000000000000,
		0b10000001000000100000010000000000000000000,
		0b100000010000001000000100000000000000000000,
		0b1000010000100001000,
		0b10000100001000010000,
		0b100001000010000100000,
		0b1000010000100001000000000,
		0b10000100001000010000000000,
		0b100001000010000100000000000,
		0b1000010000100001000000000000000,
		0b10000100001000010000000000000000,
		0b100001000010000100000000000000000,
		0b1000010000100001000000000000000000000,
		0b10000100001000010000000000000000000000,
		0b100001000010000100000000000000000000000,
	}
	for _, move := range winning_moves {
		if move&board.XPlayer == move {
			return true, 1
		}

		if move&board.OPlayer == move {
			return true, -1
		}
	}

	if board.OPlayer|board.XPlayer == FULL_BOARD {
		return true, 0
	}

	return false, 0
}

func (board Board) Player() uint8 {
	// 0 is X
	// 1 is O
	return board.TurnNumber % 2
}

func (board Board) Actions() []Board {
	// is_terminal, _ := board.IsTerminalWithValue()
	// if is_terminal {
	// 	return nil
	// }
	player := board.Player()
	var actions []Board
	combined_board := board.XPlayer | board.OPlayer

	var current_bit_address *uint64;
	for column_index := 0; column_index < 7; column_index++ {
		valid_action_column := (0b1000000 >> (bits.LeadingZeros64(combined_board & FULL_COLUMN)-58)) & FULL_COLUMN
		combined_board >>= 6
		if valid_action_column == 0 {
			continue
		}
		action := board
		action.TurnNumber += 1
		if player == X_MAX_PLAYER {
			current_bit_address = &action.XPlayer
		} else {
			current_bit_address = &action.OPlayer
		}
		*current_bit_address |= uint64(valid_action_column) << (6*uint64(column_index))
		actions = append(actions, action)
	}

	return actions
}

func (board Board) GetActionsForNumberOfCPU() ([][]Board, int) {
	number_of_cpus := runtime.NumCPU()
	// number_of_cpus := 8

	current_turn_number := board.TurnNumber
	var actions_queue []Board
	actions_queue = append(actions_queue, board.Actions()...)
	for len(actions_queue) < number_of_cpus {
		current_turn_number += 1
		for actions_queue[0].TurnNumber == current_turn_number {
			actions_queue = append(actions_queue, actions_queue[0].Actions()...)
			actions_queue = actions_queue[1:]
		}
	}

	var cpu_actions_list = make([][]Board, number_of_cpus)
	for index, action := range actions_queue {
		cpu_actions_list[index % number_of_cpus] = append(cpu_actions_list[index % number_of_cpus], action)
	}

	return cpu_actions_list, len(cpu_actions_list)
}
//                                  X win draw O win
func (board Board) Minimax(searched_games *uint64) (int8, int, int, int) {
	is_terminal, board_value := board.IsTerminalWithValue()
	player := board.Player()
	if is_terminal {
		*searched_games += 1
		if *searched_games % 1000000000 == 0 {
			fmt.Printf("searched_games: %v\n", *searched_games)
			fmt.Println(time.Now())
		}
		switch (board_value) {
			case 1: return board_value, 1, 0, 0
			case 0: return board_value, 0, 1, 0
			case -1: return board_value, 0, 0, 1
		}
	}

	x_win_count := 0
	draw_count := 0
	o_win_count := 0

	if player == X_MAX_PLAYER {
		board_value = -0b1111111
		
		for _, action := range board.Actions() {
			action_value, x_wins, draws, o_wins := action.Minimax(searched_games)

			x_win_count += int(x_wins)
			draw_count += int(draws)
			o_win_count += int(o_wins)
			board_value = max(board_value, action_value)
		}
		return board_value, x_win_count, draw_count, o_win_count
	}

	if player == O_MIN_PLAYER {
		board_value = 0b1111111
		for _, action := range board.Actions() {
			action_value, x_wins, draws, o_wins := action.Minimax(searched_games)
			x_win_count += int(x_wins)
			draw_count += int(draws)
			o_win_count += int(o_wins)
			board_value = min(board_value, action_value)
		}
		return board_value, x_win_count, draw_count, o_win_count
	}

	return 0, 0, 0, 0
}

func (board Board) IsEqual (other Board) bool {
	return board.OPlayer == other.OPlayer && board.XPlayer == other.XPlayer && board.TurnNumber == other.TurnNumber
}

func (board Board) MinimaxPrecalculatedActions(precalculated_actions []EvaluatedAction) (int8, int, int, int) {
	is_terminal, board_value := board.IsTerminalWithValue()
	player := board.Player()
	if is_terminal {
		switch (board_value) {
			case 1: return board_value, 1, 0, 0
			case 0: return board_value, 0, 1, 0
			case -1: return board_value, 0, 0, 1
		}
	}

	for _, action := range precalculated_actions {
		if board.IsEqual(action.Action) {
			return action.ActionValue, action.XWinCount, action.Draws, action.OWinCount
		}
	}

	if player == X_MAX_PLAYER {
		board_value = -0b1111111
		x_win_count := 0
		draw_count := 0
		o_win_count := 0
		for _, action := range board.Actions() {
			action_value, x_wins, draws, o_wins := action.MinimaxPrecalculatedActions(precalculated_actions)

			x_win_count += int(x_wins)
			draw_count += int(draws)
			o_win_count += int(o_wins)
			board_value = max(board_value, action_value)
		}
		return board_value, x_win_count, draw_count, o_win_count
	}

	if player == O_MIN_PLAYER {
		board_value = 0b1111111
		x_win_count := 0
		draw_count := 0
		o_win_count := 0
		for _, action := range board.Actions() {
			action_value, x_wins, draws, o_wins := action.MinimaxPrecalculatedActions(precalculated_actions)
			x_win_count += int(x_wins)
			draw_count += int(draws)
			o_win_count += int(o_wins)
			board_value = min(board_value, action_value)
		}
		return board_value, x_win_count, draw_count, o_win_count
	}
	
	return 0, 0, 0, 0
}



func (board Board) ToStringRaw() (*string, error) {
	var sb strings.Builder
	var top_row uint64 = 0b100000100000100000100000100000100000100000
	for row_index := 0; row_index < 6; row_index++ {
		x_row := (board.XPlayer & (top_row >> uint64(row_index))) >> (5-uint64(row_index))
		o_row := (board.OPlayer & (top_row >> uint64(row_index))) >> (5-uint64(row_index))
		for column_index := 6; column_index >= 0; column_index-- {
			if (x_row >> (column_index*6)) & FULL_COLUMN == 1 {
				sb.WriteRune('#')
				continue
			}
			if (o_row >> (column_index*6)) & FULL_COLUMN == 1 {
				sb.WriteRune('/')
				continue
			}
			sb.WriteRune('-')
		}
	}
	output_string := sb.String()
	return &output_string, nil
}

func (board Board) ToString() (*string, error) {

	output_string, _ := board.ToStringRaw()
	var output_sb strings.Builder
	for index, char := range []rune(*output_string) {
		output_sb.WriteRune(char)
		output_sb.WriteRune(' ')
		if (index + 1) % 7 == 0{
			output_sb.WriteRune('\n')
		}
	}
	board_string := output_sb.String()
	return &board_string, nil
}

func (board Board) Print() {
	output, err := board.ToString()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Print(*output)
}

type EvaluatedAction struct {
	Action Board
	ActionValue int8
	Score float32
	PositionNumber int
	XWinCount int
	Draws int
	OWinCount int
}

func (board Board) FindBestNextMove(searched_games *uint64) EvaluatedAction {
	player := board.Player()
	actions := board.Actions()
	evaluated_actions := make([]EvaluatedAction,len(actions))
	for index, action := range actions {
		action_value, x_wins, draws, o_wins := action.Minimax(searched_games)
		evaluated_actions[index] = EvaluatedAction{
			Action: action,
			ActionValue: action_value,
			Score: float32(x_wins*1+o_wins*-1)/float32(x_wins+draws+o_wins),
			XWinCount: x_wins,
			Draws: draws,
			OWinCount: o_wins,
		}
	}

	if player == X_MAX_PLAYER {
		
		sort.Slice(evaluated_actions, func(i, j int) bool {
			return evaluated_actions[i].Score > evaluated_actions[j].Score
		})
		sort.Slice(evaluated_actions, func(i, j int) bool {
			return evaluated_actions[i].ActionValue > evaluated_actions[j].ActionValue
		})
		return evaluated_actions[0]
	}

	if player == X_MAX_PLAYER {
		
		sort.Slice(evaluated_actions, func(i, j int) bool {
			return evaluated_actions[i].Score > evaluated_actions[j].Score
		})
		sort.Slice(evaluated_actions, func(i, j int) bool {
			return evaluated_actions[i].ActionValue > evaluated_actions[j].ActionValue
		})
		return evaluated_actions[0]
	}

	return EvaluatedAction{}
}

func (board Board) FindBestNextMovePrecalculatedActions(precalculated_actions []EvaluatedAction) EvaluatedAction {
	player := board.Player()
	actions := board.Actions()
	evaluated_actions := make([]EvaluatedAction,len(actions))
	var wg sync.WaitGroup
	for index, action := range actions {
		wg.Add(1)
		go func() {
			defer wg.Done()
			action_value, x_wins, draws, o_wins := action.MinimaxPrecalculatedActions(precalculated_actions)
			evaluated_actions[index] = EvaluatedAction{
				Action: action,
				ActionValue: action_value,
				Score: float32(x_wins*1+o_wins*-1)/float32(x_wins+draws+o_wins),
				XWinCount: x_wins,
				Draws: draws,
				OWinCount: o_wins,
			}
		}()
	}
	wg.Wait()

	if player == X_MAX_PLAYER {
		
		sort.Slice(evaluated_actions, func(i, j int) bool {
			return evaluated_actions[i].Score > evaluated_actions[j].Score
		})
		sort.Slice(evaluated_actions, func(i, j int) bool {
			return evaluated_actions[i].ActionValue > evaluated_actions[j].ActionValue
		})
		return evaluated_actions[0]
	}

	if player == X_MAX_PLAYER {
		
		sort.Slice(evaluated_actions, func(i, j int) bool {
			return evaluated_actions[i].Score > evaluated_actions[j].Score
		})
		sort.Slice(evaluated_actions, func(i, j int) bool {
			return evaluated_actions[i].ActionValue > evaluated_actions[j].ActionValue
		})
		return evaluated_actions[0]
	}

	return EvaluatedAction{}
}

func (board Board) GetNextMoveFromInput() Board {
	var input int
	player := board.Player()
	combined_board := board.XPlayer | board.OPlayer
	var valid_action_column uint64
	for {
		_, err := fmt.Scan(&input)
		if (err != nil) || !(input >= 1 && input <= 7) {
			fmt.Println("input a number 1 to 7")
			continue
		}
		valid_action_column = (0b1000000 >> (bits.LeadingZeros64(((combined_board >> (6*uint64(7-input))) & FULL_COLUMN))-58) & FULL_COLUMN)
		if valid_action_column != 0 {
			break
		}
		fmt.Println("cant place it in that column")
	}
	
	var current_bit_address *uint64
	board.TurnNumber += 1
	if player == X_MAX_PLAYER {
		current_bit_address = &board.XPlayer
	} else {
		current_bit_address = &board.OPlayer
	}
	*current_bit_address |= uint64(valid_action_column) << (6*uint64(7-input))
	return board
}

func PrintActionList(actions []EvaluatedAction) {

	max_boards_in_a_row := 7
	row_count := int(math.Ceil(float64(len(actions))/float64(max_boards_in_a_row)))
	lines := make([]string, row_count*11)
	for board_row_index := 0; board_row_index < row_count; board_row_index++  {
		for _, action := range actions[max_boards_in_a_row*board_row_index:min(max_boards_in_a_row*(board_row_index+1),len(actions))] {
			board_string, _ := action.Action.ToStringRaw()
			b := *board_string
			for row_index := 0; row_index < 6; row_index++ {
				ro := 7*row_index
				lines[row_index+7*board_row_index] += fmt.Sprintf("| %c %c %c %c %c %c %c ", b[0+ro], b[1+ro], b[2+ro], b[3+ro], b[4+ro], b[5+ro], b[6+ro])
			}
			lines[6+11*board_row_index] += fmt.Sprintf("|%2d       %5.2f ", action.ActionValue, action.Score)
			lines[7+11*board_row_index] += fmt.Sprintf("|X: %12d", action.XWinCount)
			lines[8+11*board_row_index] += fmt.Sprintf("|D: %12d", action.Draws)
			lines[9+11*board_row_index] += fmt.Sprintf("|O: %12d", action.OWinCount)
			lines[10+11*board_row_index] += "|---------------"
		}
		
	}
	for _, line := range lines {
		fmt.Printf("%v|\n", line)
	}
}

func PrintBoardList(actions []Board) {
	
	max_boards_in_a_row := 7
	row_count := int(math.Ceil(float64(len(actions))/float64(max_boards_in_a_row)))
	lines := make([]string, row_count*7)
	for board_row_index := 0; board_row_index < row_count; board_row_index++  {
		for _, action := range actions[max_boards_in_a_row*board_row_index:min(max_boards_in_a_row*(board_row_index+1),len(actions))] {
			board_string, _ := action.ToStringRaw()
			b := *board_string
			for row_index := 0; row_index < 6; row_index++ {
				ro := 7*row_index
				lines[row_index+7*board_row_index] += fmt.Sprintf("| %c %c %c %c %c %c %c ", b[0+ro], b[1+ro], b[2+ro], b[3+ro], b[4+ro], b[5+ro], b[6+ro])
			}
			lines[6+7*board_row_index] += "|---------------"
		}
		
	}

	for _, line := range lines {
		fmt.Printf("%v|\n", line)
	}
}