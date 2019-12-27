package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pascalbou/maze/lib"
)

func InitHandler(w http.ResponseWriter, r *http.Request) {
	type initRes struct {
		Name        string
		CurrentRoom int
		Cooldown    float32
		North       int `json:",omitempty"`
		East        int `json:",omitempty"`
		South       int `json:",omitempty"`
		West        int `json:",omitempty"`
	}

	type initReq struct {
		Token string
	}

	var res initRes
	var token initReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&token)
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

	q := `
	SELECT account.name, account.current_room, account.cooldown, room.north, room.east, room.south, room.west FROM account INNER JOIN room ON (account.current_room = room.room_id) WHERE account.token=$1;
	`
	rows, err := db.Query(q, token.Token)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var cooldown int64

	for rows.Next() {
		if err := rows.Scan(&res.Name, &res.CurrentRoom, &cooldown, &res.North, &res.East, &res.South, &res.West); err != nil {
			log.Fatal(err)
		}
	}

	cooldown = lib.GetCooldown(cooldown)
	res.Cooldown = float32(cooldown) / 1000

	response, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}
