package api

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pascalbou/maze/lib"
)

func CreateDB() {
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
	drop table if exists account;
	create table account (
		id SERIAL PRIMARY KEY,
		name VARCHAR(16) UNIQUE NOT NULL,
		token VARCHAR (64) UNIQUE NOT NULL,
		current_room INT NOT NULL,
		cooldown BIGINT
	);
	`

	_, err = db.Exec(q)
	if err != nil {
		log.Fatal(err)
	}

}
