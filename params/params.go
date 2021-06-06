package params

// Game holds information about the current game.
type Game struct {
	ID      string `json:"id"`
	Timeout int32  `json:"timeout"`
}

// Coord holds x/y coordinates on the board.
type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// OffBoard reports whether the coordinate is off the given board.
func (c Coord) OffBoard(board Board) bool {
	return c.X < 0 || c.Y < 0 || c.X >= board.Width || c.Y >= board.Height
}

// Distance returns the cells distance between the coordinates.
func (c Coord) Distance(other Coord) int {
	return abs(c.X-other.X) + abs(c.Y-other.Y)
}

// CloserFood returns the coordinates and the distance of the food closer to
// this coordinate, but not corresponding to the coordinate. A distance of 0 is
// returned if there is no food in the board.
func (c Coord) CloserFood(board Board) (food Coord, distance int) {
	for _, f := range board.Food {
		if f == c {
			continue
		}
		if v := c.Distance(f); v < distance || distance == 0 {
			food = f
			distance = v
		}
	}
	return food, distance
}

// CloserSnake returns the coordinates and the distance of the snake head closer to
// this coordinate, but not corresponding to the coordinate. A distance of 0 is
// returned if there are no snakes in the board.
func (c Coord) CloserSnake(board Board) (snake Battlesnake, distance int) {
	for _, s := range board.Snakes {
		if s.Head == c {
			continue
		}
		if v := c.Distance(s.Head); v < distance || distance == 0 {
			snake = s
			distance = v
		}
	}
	return snake, distance
}

// abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// OverFood reports whether the coordinate is over food.
func (c Coord) OverFood(board Board) bool {
	for _, coord := range board.Food {
		if coord == c {
			return true
		}
	}
	return false
}

// Battlesnake represents a snake, with its health and coordinates.
type Battlesnake struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Health int32   `json:"health"`
	Body   []Coord `json:"body"`
	Head   Coord   `json:"head"`
	Length int32   `json:"length"`
	Shout  string  `json:"shout"`
}

// Board holds the board state at a given turn.
type Board struct {
	Height int           `json:"height"`
	Width  int           `json:"width"`
	Food   []Coord       `json:"food"`
	Snakes []Battlesnake `json:"snakes"`
}

// GameRequest holds the overall state at a given turn.
type GameRequest struct {
	Game  Game        `json:"game"`
	Turn  int         `json:"turn"`
	Board Board       `json:"board"`
	You   Battlesnake `json:"you"`
}

// MoveResponse is the response to move requests, representing the direction a
// snake should move to.
type MoveResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout,omitempty"`
}

// BattlesnakeInfoResponse provides information about a snake.
type BattlesnakeInfoResponse struct {
	APIVersion string `json:"apiversion"`
	Author     string `json:"author"`
	Color      string `json:"color"`
	Head       string `json:"head"`
	Tail       string `json:"tail"`
}
