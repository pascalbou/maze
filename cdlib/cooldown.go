package cdlib

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/pascalbou/maze/lib"
)

func CreateCooldown(seconds time.Duration) int64 {
	return time.Now().Add(time.Second*seconds).UnixNano() / int64(time.Millisecond)
}

func GetCooldown(cooldown int64) int64 {
	result := cooldown - time.Now().UnixNano()/int64(time.Millisecond)
	return result
}

func CanAct(token string) float32 {

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
	SELECT account.cooldown FROM account WHERE account.token=$1;
	`
	rows, err := db.Query(sqlStatement, token)
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
	var cooldownTime time.Time

	cooldown = GetCooldown(cooldown)
	if cooldown > 0 {
		// adds 15s then reconvert to ms
		cooldownTime = time.Unix(0, cooldown*int64(time.Millisecond)).Add(time.Second * 15)
		fmt.Println(cooldownTime)
		cooldown = cooldownTime.UnixNano() / int64(time.Millisecond)
		fmt.Println(cooldown)
		res.Cooldown = float32(cooldown) / 1000

		sqlStatement := `
		UPDATE account SET cooldown = $2 WHERE account.token=$1;
		`
		_, err = db.Exec(sqlStatement, token, cooldown)
		if err != nil {
			log.Fatal(err)
		}

		return res.Cooldown
	}
	return 0

}
