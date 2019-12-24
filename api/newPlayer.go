package api

import (
	"database/sql"
	"log"
)

func NewPlayer(name, token string) {
	connStr := "user=postgres password=test1234 dbname=maze"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	q := `
	insert into account (name, token)
	values
		($1, $2)
	`

	_, err1 := db.Exec(q, name, token)
	if err1 != nil {
		panic(err1)
	}

}