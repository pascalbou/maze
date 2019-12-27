package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/pascalbou/maze/lib"
)

func NewPlayerHandler(w http.ResponseWriter, r *http.Request) {

	type npReq struct {
		Name  string
	}

	type npRes struct {
		Token string
	}

	var req npReq
	var res npRes
	res.Token = lib.TokenGenerator()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		log.Fatal(err)
	}

	connStr := "user=postgres password=test1234 dbname=maze"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	sqlStatement := `
	insert into account (name, token, current_room)
	values
	($1, $2, 1)
	`
	
	_, errQ := db.Exec(sqlStatement, req.Name, res.Token)
	if errQ != nil {
		log.Fatal(errQ)
	}
	
	response, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
