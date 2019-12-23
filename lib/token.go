package lib

import (
	"crypto/rand"
	"fmt"
)

func TokenGenerator() string {
	result := make([]byte, 32)
	rand.Read(result)
	return fmt.Sprintf("%x", result)
}

// func main() {
// 	token := tokenGenerator()
// 	fmt.Println(token)
// }
