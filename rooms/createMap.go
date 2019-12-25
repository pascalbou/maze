package rooms

import (
	"database/sql"
	"log"
)

type room struct {
	id          uint16
	title       string
	description string
	exitNorth   *room
	exitEast    *room
	exitSouth   *room
	exitWest    *room
	x           int8
	y           int8
}

func CreateRooms() {
	var rooms [500]room
	// var wholeMap [101][101]uint8

	var i uint16
	for i = 0; i < 3; i++ {
		rooms[i] = room{
			id:          i,
			title:       "title",
			description: "description",
		}
	}

	connStr := "user=postgres password=test1234 dbname=maze"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var sqlStatement [16]string

	sqlStatement[0] = `
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

	sqlStatement[1] = `
	insert into room (title, description, exit_north, exit_south)
	values
		($1, $2, true, false),
		($1, $2, true, true),
		($1, $2, false, true)
	`

	_, err1 := db.Exec(sqlStatement[0])
	if err1 != nil {
		log.Fatal(err1)
	}

	_, err2 := db.Exec(sqlStatement[1], "title", "description")
	if err2 != nil {
		log.Fatal(err2)
	}
}
