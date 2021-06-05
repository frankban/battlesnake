package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/frankban/battlesnake/api"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	http.HandleFunc("/", api.HandleIndex)
	http.HandleFunc("/start", api.HandleStart)
	http.HandleFunc("/move", api.HandleMove)
	http.HandleFunc("/end", api.HandleEnd)

	fmt.Printf("starting battlesnake server at http://0.0.0.0:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
