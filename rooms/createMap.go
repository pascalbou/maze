package rooms

import (
	"database/sql"
	"log"
)

type room struct {
	id          int
	title       string
	description string
	// players     []string
	// items       []string
	exits []string
}

func CreateRooms() {
	var rooms [3]room

	for i := 0; i < 3; i++ {
		rooms[i] = room{
			id:          i,
			title:       "title",
			description: "description",
		}
	}

	rooms[0].exits = append(rooms[0].exits, "n")
	rooms[1].exits = append(rooms[1].exits, "n", "s")
	rooms[2].exits = append(rooms[2].exits, "s")

	connStr := "user=postgres password=test1234 dbname=maze"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStatement0 := `
	drop table if exists room;
	create table room (
		id SERIAL,
		title VARCHAR(64),
		description VARCHAR(256),
		exit_north boolean,
		exit_east boolean,
		exit_south boolean,
		exit_west boolean
	);
	`

	sqlStatement1 := `
	insert into room (title, description, exit_north, exit_south)
	values
		($1, $2, true, false),
		($1, $2, true, true),
		($1, $2, false, true)
	`

	_, err1 := db.Exec(sqlStatement0)
	if err1 != nil {
		log.Fatal(err1)
	}

	_, err2 := db.Exec(sqlStatement1, "title", "description")
	if err2 != nil {
		log.Fatal(err2)
	}
}


