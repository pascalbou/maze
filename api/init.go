package api

import (
	"database/sql"
	"log"
)

func Init(token string) (string, int) {
	connStr := "user=postgres password=test1234 dbname=maze"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	q := `
	SELECT name, current_room FROM account WHERE token=$1
	`

	rows, err1 := db.Query(q, token)
	if err1 != nil {
		panic(err1)
	}
	defer rows.Close()

	var name string
	var room int
	for rows.Next() {
		if err := rows.Scan(&name, &room); err != nil {
			panic(err)
		}
	}

	// fmt.Println(name, room)

	return name, room

}
