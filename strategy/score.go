package strategy

import (
	"fmt"
	"math/rand"

	"github.com/frankban/battlesnake/params"
)

// moveByScore is a rudimentary strategy based on figuring out a score for every
// possible move, by looking at free cells, food presence and enemies.
func moveByScore(state *params.GameRequest, ds []Direction) Direction {
	rand.Shuffle(len(ds), func(i, j int) { ds[i], ds[j] = ds[j], ds[i] })
	var result scoreResult
	results := make(chan scoreResult)

	for _, d := range ds {
		d := d
		go func() {
			l := &logger{
				prefix: "    ",
			}
			score := getScore(state, d, l)
			results <- scoreResult{
				direction: d,
				score:     score,
				logs:      l.String(),
			}
		}()
	}

	for range ds {
		sc := <-results
		fmt.Printf("  score going %s: %d\n%s", sc.direction, sc.score, sc.logs)
		if sc.score > result.score {
			result = sc
		}
	}

	return result.direction
}

func getScore(state *params.GameRequest, d Direction, l *logger) (score int) {
	s := nextSnake(state.You, state.Board, d)
	board := nextBoard(s, state.Board)

	// Calculate free cells available after this move.
	score = freeCellsFrom(board, s.Head)
	l.Log("%d from free cells\n", score)

	if s.Health < 50 && s.Head.OverFood(board) {
		// The head is over food.
		sc := 1 + int(50/s.Health)
		l.Log("%d because over food\n", sc)
		score += sc
	} else if s.Health < 10 {
		// How closer the snake gets to food if it is starving?
		if _, distance := s.Head.CloserFood(board); distance != 0 {
			sc := (state.Board.Height+state.Board.Width)/2 - distance
			l.Log("%d for getting closer to food\n", sc)
			score += sc
		}
	}

	// Is the head close to another snake's head after this move, and can
	// they collide?
	if snake, distance := s.Head.CloserSnake(board); !isEven(distance) {
		var sc int
		if state.You.Length > snake.Length {
			if distance == 1 {
				sc = 1
				l.Log("%d for getting closer to shorter snake\n", sc)
			}
		} else {
			sc = -((state.Board.Height+state.Board.Width)/2 - distance)
			l.Log("%d for getting closer to longer snake\n", sc)
		}
		score += sc
	}
	return score
}

type scoreResult struct {
	direction Direction
	score     int
	logs      string
}
