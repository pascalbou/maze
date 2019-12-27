package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pascalbou/maze/lib"
)

func MoveHandler(w http.ResponseWriter, r *http.Request) {
	type moveReq struct {
		Token     string
		Direction string
	}
	
	type moveRes struct {
		PreviousRoom int
		CurrentRoom    int
		Message     string
	}

	var req moveReq
	var res moveRes
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

	q := `
	SELECT current_room FROM account WHERE account.token=$1;
	`
	rows, err := db.Query(q, req.Token)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&res.PreviousRoom); err != nil {
			log.Fatal(err)
		}
	}

	q2 := fmt.Sprintf(`
	SELECT %s FROM room WHERE room.room_id=$1;
	`, req.Direction)

	rows, err = db.Query(q2, res.PreviousRoom)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&res.CurrentRoom); err != nil {
			log.Fatal(err)
		}
	}

	if res.CurrentRoom != 0 {
		res.Message = "You moved " + req.Direction
		q := `
		UPDATE account SET current_room = $2 WHERE account.token=$1;
		`

		_, err := db.Exec(q, req.Token, res.CurrentRoom)
		if err != nil {
			log.Fatal(err)
		} 
	} else {
		res.Message = "You cannot move " + req.Direction
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
