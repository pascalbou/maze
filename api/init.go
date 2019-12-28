package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pascalbou/maze/cdlib"
	"github.com/pascalbou/maze/lib"
)

func InitHandler(w http.ResponseWriter, r *http.Request) {
	type initRes struct {
		Name        string
		CurrentRoom int
		Cooldown    float32
		// North       int `json:",omitempty"`
		// East        int `json:",omitempty"`
		// South       int `json:",omitempty"`
		// West        int `json:",omitempty"`
		North       int `json:"-"`
		East        int `json:"-"`
		South       int `json:"-"`
		West        int `json:"-"`
		Exits []string
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

	var cooldownDB int64

	for rows.Next() {
		if err := rows.Scan(&res.Name, &res.CurrentRoom, &cooldownDB, &res.North, &res.East, &res.South, &res.West); err != nil {
			log.Fatal(err)
		}
	}

	cooldown := cdlib.GetCooldown(cooldownDB)
	if cooldown < 0 {
		cooldown = 1000
	}
	res.Cooldown = float32(cooldown) / 1000

	// var list []int
	// list = append(list, "North", "East", "South", "West")

	// for _, element := list {
	// 	if res[string(element)] != 0 {
	// 		res.Exits = append(res.Exits, string(element))
	// 	}
	// }

	if res.North != 0 {
		res.Exits = append(res.Exits, "north")
	}
	if res.East != 0 {
		res.Exits = append(res.Exits, "east")
	}
	if res.South != 0 {
		res.Exits = append(res.Exits, "south")
	}
	if res.West != 0 {
		res.Exits = append(res.Exits, "west")
	}

	// fmt.Println(res.Exits)

	response, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}
