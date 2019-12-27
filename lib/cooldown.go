package lib

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pascalbou/maze/lib"
)

func AddCooldown(seconds time.Duration) int64 {
	return time.Now().Add(time.Second*seconds).UnixNano() / int64(time.Millisecond)
}

func GetCooldown(cooldown int64) int64 {
	result := cooldown - time.Now().UnixNano()/int64(time.Millisecond)
	if result < 0 {
		result = 1000
	}
	return result
}

func CanAct(token string) bool {

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
	SELECT account.cooldown FROM account WHERE account.token=$1;
	`
	rows, err := db.Query(q, token)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var cooldown int64

	for rows.Next() {
		if err := rows.Scan(&cooldown); err != nil {
			log.Fatal(err)
		}
	}

	type caBody struct {
		Message  string
		Cooldown float32
	}
	var res caBody

	cooldown -= time.Now().UnixNano() / int64(time.Millisecond)
	if cooldown > 0 {
		// not working, need to add current cooldown + 15
		cooldown = AddCooldown(15)
		res.Message = "You acted before your cooldown finished. Penalty +15s."
		res.Cooldown = float32(cooldown) / 1000

		sqlStatement := `
		UPDATE account SET cooldown = $2 WHERE account.token=$1;
		values
		($1, $2, 1, $3)
		`
		cooldown := lib.AddCooldown(30)

		_, err = db.Exec(sqlStatement, token, cooldown)
		if err != nil {
			log.Fatal(err)
		}

		response, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)

		return false
	}
	return true

}
