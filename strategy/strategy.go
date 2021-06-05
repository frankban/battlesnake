package strategy

import (
	"fmt"
	"math/rand"

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
	// start := time.Now()

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

	// // Then refine the selection, but without running out of time.
	// badMoves := make(chan Direction, 1)
	// badMoves := refine()
	// d := state.Game.Timeout*time.Millisecond - time.Since(start).Milliseconds() - 50*time.Millisecond
	// for {
	// 	select {
	// 	case badMove := <-badMoves:

	// 		fmt.Println(res)
	// 	case <-time.After(d):
	// 		fmt.Println("running out of time")
	// 		break
	// 	}
	// }

	return ds[rand.Intn(len(ds))]
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
