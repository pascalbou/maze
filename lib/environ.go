package lib

import (
	"log"

	"github.com/joho/godotenv"
)

func GetEnviron() map[string]string {
	var myEnv map[string]string
	myEnv, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// dbName := myEnv["DB_Name"]

	return myEnv

}
