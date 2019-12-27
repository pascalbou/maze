package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type moveReq struct {
	Token     string
	Direction string
}

type moveRes struct {
	CurrentRoom int
	NextRoom    int
	Message     string
}

func MoveHandler(w http.ResponseWriter, r *http.Request) {
	var req moveReq
	var res moveRes
	decoder := json.NewDecoder(r.Body)
	errD := decoder.Decode(&req)
	if errD != nil {
		panic(errD)
	}

	connStr := "user=postgres password=test1234 dbname=maze"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// roomColumn := "room." + req.Direction

	q := `
	SELECT current_room FROM account WHERE account.token=$1;
	`

	rows, err1 := db.Query(q, req.Token)
	if err1 != nil {
		panic(err1)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&res.CurrentRoom); err != nil {
			panic(err)
		}
	}

	q2 := fmt.Sprintf(`
	SELECT %s FROM room WHERE room.room_id=$1;
	`, req.Direction)

	rows2, err2 := db.Query(q2, res.CurrentRoom)
	if err2 != nil {
		panic(err2)
	}
	defer rows2.Close()

	for rows2.Next() {
		if err := rows2.Scan(&res.NextRoom); err != nil {
			panic(err)
		}
	}

	if res.NextRoom != 0 {
		res.Message = "You moved " + req.Direction
		q := `
		UPDATE account SET current_room = $2 WHERE account.token=$1;
		`

		_, err1 := db.Exec(q, req.Token, res.NextRoom)
		if err1 != nil {
			panic(err1)
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
