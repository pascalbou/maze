package main

import (
	"crypto/rand"
	"fmt"
)

func tokenGenerator() string {
	result := make([]byte, 32)
	rand.Read(result)
	return fmt.Sprintf("%x", result)
}

func main() {
	token := tokenGenerator()
	fmt.Println(token)
}
