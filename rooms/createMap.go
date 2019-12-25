package rooms

import (
	"fmt"
	"math/rand"
	"time"
)

type room struct {
	id          int
	title       string
	description string
	// north   room
	// east    room
	// south   room
	// west	room
	exits map[string]room
	x     int
	y     int
}

func createOneRoom(id, x, y int) room {
	return room{
		id: id,
		x:  x,
		y:  y,
	}
}

// func createTenRandomRooms(id int) {

// }

func getDirection() string {
	var directions []string
	directions = append(directions, "north", "east", "south", "west")
	return directions[rand.Intn(len(directions))]
}

func getOppositeDirectopn(direction string) string {
	var opposite string
	switch direction {
	case "north":
		opposite = "south"
	case "east":
		opposite = "west"
	case "south":
		opposite = "north"
	case "west":
		opposite = "east"
	}
	return opposite
}

func getCoordinates(direction string, x, y int) (int, int) {
	switch direction {
	case "north":
		y++
	case "east":
		x++
	case "south":
		y--
	case "west":
		x--
	}
	return x, y
}

func conditionsRoom(existingRoom, x, y int) bool {
	if existingRoom != 0 {
		return false
	}
	if x < -50 || x > 50 {
		return false
	}
	if y < -50 || y > 50 {
		return false
	}
	return true
}

func CreateRooms() {
	rand.Seed(time.Now().Unix())
	var rooms [500]room
	wholeMaze := make(map[string]string)

	var x, y int
	x = 0
	y = 0

	rooms[0] = createOneRoom(0, x, y)
	// countRooms := 1
	wholeMaze["0,0"] = "0"
	previousRoom := rooms[0]
	previousRoom.exits = make(map[string]room)
	var currentRoom room

	for i := 0; i < 10; i++ {
		randomDirection := getDirection()
		oppositeDirection := getOppositeDirectopn(randomDirection)
		x, y = getCoordinates(randomDirection, x, y)
		currentRoom = createOneRoom(i, x, y)
		keyWholeMaze := fmt.Sprintf("%d,%d", x, y)

		fmt.Println(randomDirection)
		fmt.Println(oppositeDirection)
		fmt.Println(currentRoom)
		fmt.Println(keyWholeMaze)
		fmt.Println(len(wholeMaze[keyWholeMaze]))
		fmt.Println(x)
		fmt.Println(y)

		if conditionsRoom(len(wholeMaze[keyWholeMaze]), x, y) {
			currentRoom.exits = make(map[string]room)
			previousRoom.exits[randomDirection] = currentRoom
			currentRoom.exits[oppositeDirection] = previousRoom
			wholeMaze[keyWholeMaze] = string(i)
			previousRoom = currentRoom
		} else {
			i--
		}
	}

	// connStr := "user=postgres password=test1234 dbname=maze"
	// db, err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// var sqlStatement [16]string

	// sqlStatement[0] = `
	// drop table if exists room;
	// create table room (
	// 	id SERIAL,
	// 	title VARCHAR(64),
	// 	description VARCHAR(256),
	// 	exit_north boolean,
	// 	exit_east boolean,
	// 	exit_south boolean,
	// 	exit_west boolean
	// );
	// `

	// sqlStatement[1] = `
	// insert into room (title, description, exit_north, exit_south)
	// values
	// 	($1, $2, true, false),
	// 	($1, $2, true, true),
	// 	($1, $2, false, true)
	// `

	// _, err1 := db.Exec(sqlStatement[0])
	// if err1 != nil {
	// 	log.Fatal(err1)
	// }

	// _, err2 := db.Exec(sqlStatement[1], "title", "description")
	// if err2 != nil {
	// 	log.Fatal(err2)
	// }
}
