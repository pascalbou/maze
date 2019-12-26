package rooms

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type room struct {
	id int
	// north   room
	// east    room
	// south   room
	// west	room
	exits map[string]int
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
	rooms := make(map[int]room)
	wholeMaze := make(map[string]string)

	var x, y int
	x = 0
	y = 0

	rooms[1] = createOneRoom(1, x, y)
	countRooms := 1
	wholeMaze["0,0"] = "1"

	totalRooms := 20

	for countRooms < totalRooms {
		remainingRooms := totalRooms - countRooms
		if remainingRooms > 10 {
			remainingRooms = 10
		}
		previousRoom := rooms[rand.Intn(len(rooms))]
		if previousRoom.exits == nil {
			previousRoom.exits = make(map[string]int)
		}
		for i := 0; i < remainingRooms; i++ {
			randomDirection := getDirection()
			oppositeDirection := getOppositeDirectopn(randomDirection)
			x, y = getCoordinates(randomDirection, previousRoom.x, previousRoom.y)
			currentRoom := createOneRoom(countRooms, x, y)
			keyWholeMaze := fmt.Sprintf("%d,%d", x, y)
	
			// fmt.Println(randomDirection)
			// fmt.Println(oppositeDirection)
			// fmt.Println(currentRoom)
			// fmt.Println(keyWholeMaze)
			// fmt.Println(len(wholeMaze[keyWholeMaze]))
			// fmt.Println(x)
			// fmt.Println(y)
	
			if conditionsRoom(len(wholeMaze[keyWholeMaze]), x, y) {
				currentRoom.exits = make(map[string]int)
				previousRoom.exits[randomDirection] = currentRoom.id
				currentRoom.exits[oppositeDirection] = previousRoom.id
				wholeMaze[keyWholeMaze] = string(countRooms)
				// fmt.Println(currentRoom)
				// fmt.Println(rooms[i])
				rooms[currentRoom.id] = currentRoom
				rooms[previousRoom.id] = previousRoom
				previousRoom = currentRoom
				countRooms++
				
			} else {
				break
			}
		}
	}

	fmt.Println(countRooms)
	fmt.Println(len(rooms))
	// fmt.Println(rooms)
	// for k, v := range rooms {
	// 	fmt.Println(k, v)
	// }
	for i:=0;i<len(rooms);i++ {
		fmt.Println(i, rooms[i])
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
		room_id INT,
		north INT,
		east INT,
		south INT,
		west INT,
		x INT,
		y INT
	);
	`

	_, err1 := db.Exec(sqlStatement[0])
	if err1 != nil {
		log.Fatal(err1)
	}

	sqlStatement[1] = `
		insert into room (room_id, north, east, south, west, x, y)
		values
		`
	// sqlStatement[1] += fmt.Sprintf("\t(%d, %d, %d, %d, %d, %d, %d)\n\t\t", rooms[0].id, rooms[0].exits["north"], rooms[0].exits["east"], rooms[0].exits["south"], rooms[0].exits["west"], rooms[0].x, rooms[0].y)	

	// fmt.Println(sqlStatement[1])

	for i := 0; i < len(rooms); i++ {
		vals := fmt.Sprintf("\t(%d, %d, %d, %d, %d, %d, %d),\n\t\t", rooms[i].id, rooms[i].exits["north"], rooms[i].exits["east"], rooms[i].exits["south"], rooms[i].exits["west"], rooms[i].x, rooms[i].y)
		// fmt.Println(vals)
		sqlStatement[1] += vals
		// fmt.Println(sqlStatement[1])

	}

	sqlStatement[1] = sqlStatement[1][:len(sqlStatement[1])-4]
	// fmt.Println(sqlStatement[1])

	_, err2 := db.Exec(sqlStatement[1])
	if err2 != nil {
		log.Fatal(err2)
	}
}
