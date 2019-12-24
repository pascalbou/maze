package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pascalbou/maze/lib"
)

type body struct {
	Name  string
	Token string
}

func handler(w http.ResponseWriter, r *http.Request) {
	token := lib.TokenGenerator()
	// fmt.Println(token)

	var n body
	decoder := json.NewDecoder(r.Body)
	errD := decoder.Decode(&n)
	if errD != nil {
		panic(errD)
	}

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
	http.HandleFunc("/newplayer", handler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
