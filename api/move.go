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

func MoveHandler(w http.ResponseWriter, r *http.Request) {

	type moveReq struct {
		Token     string
		Direction string
		NextRoom int
	}

	type moveRes struct {
		PreviousRoom int
		CurrentRoom  int
		Message      string
		Cooldown     float32
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

	check := cdlib.CanAct(req.Token)

	if check != 0 {
		res.Cooldown = check
		res.Message = "You acted before your cooldown finished. Penalty +15s."

		response, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)

		return
	}

	sqlStatement := `
	SELECT current_room FROM account WHERE account.token=$1;
	`
	rows, err := db.Query(sqlStatement, req.Token)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&res.PreviousRoom); err != nil {
			log.Fatal(err)
		}
	}

	sqlStatement = fmt.Sprintf(`
	SELECT %s FROM room WHERE room.room_id=$1;
	`, req.Direction)

	rows, err = db.Query(sqlStatement, res.PreviousRoom)
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
		sqlStatement := `
		UPDATE account SET current_room = $2, cooldown = $3 WHERE account.token=$1;
		`

		cooldown := cdlib.CreateCooldown(30)
		res.Cooldown = 30
			if req.NextRoom == res.CurrentRoom {
				cooldown = cdlib.CreateCooldown(15)
				res.Cooldown /= 2
				res.Message += " . Wise explorer: -50% cooldown" 
			} else {
				cooldown = cdlib.CreateCooldown(45)
				res.Cooldown = 45
				res.Message += " but next room id is incorrect: +15s penalty"
			}

		_, err := db.Exec(sqlStatement, req.Token, res.CurrentRoom, cooldown)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		res.Message = "Cooldown penalty 60s. You cannot move " + req.Direction
		sqlStatement := `
		UPDATE account SET cooldown = $2 WHERE account.token=$1;
		`

		cooldown := cdlib.CreateCooldown(60)
		res.Cooldown = 60

		_, err := db.Exec(sqlStatement, req.Token, cooldown)
		if err != nil {
			log.Fatal(err)
		}
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
