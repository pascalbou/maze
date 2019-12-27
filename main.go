package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/pascalbou/maze/api"
	"github.com/pascalbou/maze/rooms"
)

func main() {
	// api.CreateDB()
	rooms.CreateRooms()

	http.HandleFunc("/newplayer", api.NewPlayerHandler)
	http.HandleFunc("/init", api.InitHandler)
	http.HandleFunc("/move", api.MoveHandler)

	log.Println("Listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))

}
