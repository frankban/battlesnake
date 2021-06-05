package strategy

import (
	"fmt"
	"math/rand"
	"time"

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

var directions = []Direction{up, down, left, right}

// returns the direction of battlesnake's next move.
func Move(state *params.GameRequest) Direction {
	start := time.Now()

	// First exclude moves leadign to immediate death.
	ds := possibleDirections(state.You, state.Board)
	switch len(ds) {
	case 0:
		fmt.Println("there is no tomorrow: turn left!")
		return left
	case 1:
		fmt.Println("one choice only")
		return ds[0]
	}

	// Then refine the selection, but without running out of time.
	suggestions := make(chan Direction)
	done := make(chan bool)
	go refine(state, ds, suggestions, done)
	//deadline := time.Duration(int64(state.Game.Timeout)-time.Since(start).Milliseconds()-50) * time.Millisecond
	deadline := 400 * time.Millisecond
outer:
	for {
		select {
		case d := <-suggestions:
			fmt.Printf("found suggestion while refining: %s\n", d)
			close(done)
			return d
		case <-time.After(deadline):
			fmt.Println("running out of time")
			close(done)
			break outer
		}
	}

	return ds[rand.Intn(len(ds))]
}

func refine(state *params.GameRequest, ds []Direction, suggestions chan Direction, done chan bool) {
	var suggestion Direction
	var freeCells int
	for _, d := range ds {
		s := nextSnake(state.You, state.Board, d)
		board := nextBoard(s, state.Board)
		if free := freeCellsFrom(board, s.Head); free > freeCells {
			freeCells = free
			suggestion = d
			fmt.Printf("found %d free cells going %s", freeCells, suggestion)
		}
	}
	suggestions <- suggestion
}

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

func nextSnake(s params.Battlesnake, board params.Board, d Direction) params.Battlesnake {
	s.Head = nextCoord(s.Head, d)
	s.Body = append([]params.Coord{s.Head}, s.Body...)
	for _, c := range board.Food {
		if c == s.Head {
			s.Length += 1
			return s
		}
	}
	s.Body = s.Body[:s.Length-1]
	return s
}

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
			for _, c := range s.Body[1:] {
				if next == c {
					continue outer
				}
			}
		}

		ds = append(ds, d)
	}
	return ds
}

// FreeCellsFrom returns the number of free contiguous cells in the given board
// from tne given coordinate.
func freeCellsFrom(board params.Board, c params.Coord) int {
	free := make(map[params.Coord]bool)
	taken := make(map[params.Coord]bool)
	for _, s := range board.Snakes {
		for _, c := range s.Body {
			taken[c] = true
		}
	}
	freeCellsFrom0(c, free, taken)
	return len(free)
}

func freeCellsFrom0(c params.Coord, free, taken map[params.Coord]bool) {
	for _, d := range directions {
		next := nextCoord(c, d)
		if taken[next] || free[next] {
			continue
		}
		free[next] = true
		freeCellsFrom0(next, free, taken)
	}
}
