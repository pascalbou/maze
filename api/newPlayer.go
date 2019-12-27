package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pascalbou/maze/lib"
)

func NewPlayerHandler(w http.ResponseWriter, r *http.Request) {

	type npReq struct {
		Name string
	}

	type npRes struct {
		Token    string
		Cooldown int
	}

	var req npReq
	var res npRes
	res.Token = lib.TokenGenerator()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		log.Fatal(err)
	}

	dbUser := lib.GetEnviron()["DB_USER"]
	dbPass := lib.GetEnviron()["DB_PASS"]
	dbName := lib.GetEnviron()["DB_NAME"]

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s", dbUser, dbPass, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStatement := `
	insert into account (name, token, current_room, cooldown)
	values
	($1, $2, 1, $3)
	`
	cooldown := lib.AddCooldown(30)

	_, err = db.Exec(sqlStatement, req.Name, res.Token, cooldown)
	if err != nil {
		log.Fatal(err)
	}

	res.Cooldown = 30

	response, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
