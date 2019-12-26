package api

import (
	"database/sql"
	"log"
)

func CreateDB() {

	connStr := "user=postgres password=test1234 dbname=maze"
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
		current_room INT NOT NULL
	);
	`

	_, err1 := db.Exec(q)
	if err1 != nil {
		panic(err1)
	}

}
