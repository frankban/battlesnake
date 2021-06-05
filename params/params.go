package params

type Game struct {
	ID      string `json:"id"`
	Timeout int32  `json:"timeout"`
}

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// OffBoard reports whether the coordinate is off the given board.
func (c Coord) OffBoard(b Board) bool {
	return c.X < 0 || c.Y < 0 || c.X >= b.Width || c.Y >= b.Height
}

// CloseTo reports whether the coordinate is close to the given one.
func (c Coord) CloseTo(other Coord) bool {
	return abs(c.X-other.X)+abs(c.Y-other.Y) == 1
}

// abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// OverFood reports whether the coordinate is over food.
func (c Coord) OverFood(b Board) bool {
	for _, coord := range b.Food {
		if coord == c {
			return true
		}
	}
	return false
}

type Battlesnake struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Health int32   `json:"health"`
	Body   []Coord `json:"body"`
	Head   Coord   `json:"head"`
	Length int32   `json:"length"`
	Shout  string  `json:"shout"`
}

type Board struct {
	Height int           `json:"height"`
	Width  int           `json:"width"`
	Food   []Coord       `json:"food"`
	Snakes []Battlesnake `json:"snakes"`
}

type GameRequest struct {
	Game  Game        `json:"game"`
	Turn  int         `json:"turn"`
	Board Board       `json:"board"`
	You   Battlesnake `json:"you"`
}

type MoveResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout,omitempty"`
}

type BattlesnakeInfoResponse struct {
	APIVersion string `json:"apiversion"`
	Author     string `json:"author"`
	Color      string `json:"color"`
	Head       string `json:"head"`
	Tail       string `json:"tail"`
}
