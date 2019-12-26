package main

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/pascalbou/maze/api"
	"github.com/pascalbou/maze/lib"
	"github.com/pascalbou/maze/rooms"
)

type body struct {
	Name  string
	Token string
}

func newPlayerHandler(w http.ResponseWriter, r *http.Request) {
	token := lib.TokenGenerator()
	// fmt.Println(len(token))

	var b body
	decoder := json.NewDecoder(r.Body)
	errD := decoder.Decode(&b)
	if errD != nil {
		panic(errD)
	}

	api.NewPlayer(b.Name, token)

	newBody := body{Name: b.Name, Token: token}
	response, err := json.Marshal(newBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

type initBody struct {
	Name string
	Token string
	CurrentRoom int
	// North int
	// East int
	// South int
	// West int
}

func initHandler(w http.ResponseWriter, r *http.Request) {
	var b initBody
	decoder := json.NewDecoder(r.Body)
	errD := decoder.Decode(&b)
	if errD != nil {
		panic(errD)
	}

	b.Name, b.CurrentRoom = api.Init(b.Token)
	response, err := json.Marshal(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func main() {
	// api.CreateDB()
	rooms.CreateRooms()

	http.HandleFunc("/newplayer", newPlayerHandler)
	http.HandleFunc("/init", initHandler)

	log.Println("Listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))



}
