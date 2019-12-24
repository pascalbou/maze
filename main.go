package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/pascalbou/maze/api"
	"github.com/pascalbou/maze/lib"
)

type body struct {
	Name  string
	Token string
}

func handler(w http.ResponseWriter, r *http.Request) {
	token := lib.TokenGenerator()
	// fmt.Println(len(token))

	var n body
	decoder := json.NewDecoder(r.Body)
	errD := decoder.Decode(&n)
	if errD != nil {
		panic(errD)
	}

	api.NewPlayer(n.Name, token)

	newBody := body{Name: n.Name, Token: token}
	response, err := json.Marshal(newBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func main() {
	api.CreateDB()

	http.HandleFunc("/newplayer", handler)
	log.Fatal(http.ListenAndServe(":3000", nil))

}
