package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type initBody struct {
	Name        string
	CurrentRoom int
	North       int `json:",omitempty"`
	East        int `json:",omitempty"`
	South       int `json:",omitempty"`
	West        int `json:",omitempty"`
}

type token struct {
	Token string
}

func InitHandler(w http.ResponseWriter, r *http.Request) {
	var b initBody
	var t token
	decoder := json.NewDecoder(r.Body)
	errD := decoder.Decode(&t)
	if errD != nil {
		panic(errD)
	}

	// fmt.Println(t)

	connStr := "user=postgres password=test1234 dbname=maze"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	q := `
	SELECT account.name, account.current_room, room.north, room.east, room.south, room.west FROM account INNER JOIN room ON (account.current_room = room.room_id) WHERE account.token=$1;
	`

	rows, err1 := db.Query(q, t.Token)
	if err1 != nil {
		panic(err1)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&b.Name, &b.CurrentRoom, &b.North, &b.East, &b.South, &b.West); err != nil {
			panic(err)
		}
	}

	response, err := json.Marshal(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}
