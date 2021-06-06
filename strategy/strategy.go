package strategy

import (
	"fmt"

	"github.com/frankban/battlesnake/params"
)

// Direction holds the direction of a single move.
type Direction string

const (
	up    Direction = "up"
	down  Direction = "down"
	left  Direction = "left"
	right Direction = "right"
)

// directions stores the four possible directions.
var directions = []Direction{up, down, left, right}

// Move returns the direction of battlesnake's next move.
func Move(state *params.GameRequest) Direction {
	// First exclude moves leading to immediate death.
	ds := possibleDirections(state.You, state.Board)
	switch len(ds) {
	case 0:
		fmt.Println("  there is no tomorrow: turn left!")
		return left
	case 1:
		fmt.Println("  one choice only")
		return ds[0]
	}

	// Then refine the selection.
	return moveByScore(state, ds)
}

// possibleDirections returns the possible directions that the given snake could
// take in the given board that would not lead to immediate death.
func possibleDirections(snake params.Battlesnake, board params.Board) []Direction {
	ds := make([]Direction, 0, 3)
outer:
	for _, d := range directions {
		next := nextCoord(snake.Head, d)

		// Do not go off board.
		if next.OffBoard(board) {
			continue
		}

		// Do not hit snake bodies.
		for _, s := range board.Snakes {
			for _, c := range s.Body {
				if next == c {
					continue outer
				}
			}
		}

		ds = append(ds, d)
	}
	return ds
}

// nextCoord return the coordinates resulting from starting from the given
// coordinate and moving in the given direction.
func nextCoord(c params.Coord, d Direction) params.Coord {
	switch d {
	case up:
		return params.Coord{
			X: c.X,
			Y: c.Y + 1,
		}
	case down:
		return params.Coord{
			X: c.X,
			Y: c.Y - 1,
		}
	case left:
		return params.Coord{
			X: c.X - 1,
			Y: c.Y,
		}
	default:
		return params.Coord{
			X: c.X + 1,
			Y: c.Y,
		}
	}
}

// nextSnake return the snake resulting from moving the given snake in the board
// in the given direction.
func nextSnake(s params.Battlesnake, board params.Board, d Direction) params.Battlesnake {
	s.Head = nextCoord(s.Head, d)
	s.Body = append([]params.Coord{s.Head}, s.Body...)
	if s.Head.OverFood(board) {
		s.Length += 1
		return s
	}
	s.Body = s.Body[:s.Length]
	return s
}

// nextBoard return the board that results when putting the given snake into the
// given board.
func nextBoard(s params.Battlesnake, board params.Board) params.Board {
	snakes := make([]params.Battlesnake, len(board.Snakes))
	for i, snake := range board.Snakes {
		if s.ID == snake.ID {
			snakes[i] = s
			continue
		}
		snakes[i] = snake
	}
	board.Snakes = snakes
	return board
}

// freeCellsFrom returns the number of free contiguous cells in the given board
// from tne given coordinate.
func freeCellsFrom(board params.Board, c params.Coord) int {
	free := make(map[params.Coord]bool)
	taken := make(map[params.Coord]bool)
	for _, s := range board.Snakes {
		for _, c := range s.Body {
			taken[c] = true
		}
	}
	freeCellsFrom0(c, board, free, taken)
	return len(free)
}

func freeCellsFrom0(c params.Coord, board params.Board, free, taken map[params.Coord]bool) {
	for _, d := range directions {
		next := nextCoord(c, d)
		if next.OffBoard(board) || taken[next] || free[next] {
			continue
		}
		free[next] = true
		freeCellsFrom0(next, board, free, taken)
	}
}

// isEven reports whether the given number is even.
func isEven(n int) bool {
	return n%2 == 0
}
