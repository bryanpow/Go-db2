package main

import (
	"fmt"
	"go_db/db"
)



func main() {
	// Example usage
	database := db.NewDb()
	database.Set("username", "johndoe")
	if value, exists := database.Get("username"); exists {
		fmt.Println("Retrieved:", value)
	}
	database.Delete("username")
	if _, exists := database.Get("username"); !exists {
		fmt.Println("Username key has been deleted.")
	}
}