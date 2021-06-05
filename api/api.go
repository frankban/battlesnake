package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/frankban/battlesnake/params"
	"github.com/frankban/battlesnake/strategy"
)

// HandleIndex is called when your Battlesnake is created and refreshed
// by play.battlesnake.com. BattlesnakeInfoResponse contains information about
// your Battlesnake, including what it should look like on the game board.
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	resp := params.BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "frankban",
		Color:      "#FF0000",
		Head:       "bendr",
		Tail:       "round-bum",
	}
	writeResponse(w, resp)
}

// HandleStart is called at the start of each game your Battlesnake is playing.
// The GameRequest object contains information about the game that's about to start.
// TODO: Use this function to decide how your Battlesnake is going to look on the board.
func HandleStart(w http.ResponseWriter, r *http.Request) {
	// Nothing to respond with here.
	fmt.Println("* start")
}

// HandleMove is called for each turn of each game.
// Valid responses are "up", "down", "left", or "right".
func HandleMove(w http.ResponseWriter, r *http.Request) {
	state := getState(r)
	fmt.Printf("* turn %d: start\n", state.Turn)

	resp := params.MoveResponse{
		Move: string(strategy.Move(state)),
	}

	fmt.Printf("* move %d: %s\n\n", state.Turn, resp.Move)
	writeResponse(w, resp)
}

// HandleEnd is called when a game your Battlesnake was playing has ended.
// It's purely for informational purposes, no response required.
func HandleEnd(w http.ResponseWriter, r *http.Request) {
	// Nothing to respond with here.
	fmt.Println("* end")
}

func getState(r *http.Request) (state *params.GameRequest) {
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		log.Fatal(err)
	}
	return state
}

func writeResponse(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Fatal(err)
	}
}
