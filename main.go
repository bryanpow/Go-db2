package main

import (
	"fmt"
	"log"
	"go_db/db"
)



func main() {
	// Create a new DB instance
	database := db.NewDb("my_database.db")

	// Set some data in the database
	database.Set("username", "john_doe")
	database.Set("email", "john@example.com")

	// Save the current state of the database
	err := database.Save()
	if err != nil {
		log.Fatalf("Failed to save database: %v", err)
	}

	// Create a new database instance to test loading from file
	newDatabase := db.NewDb("my_database.db")

	// Load the data into the new database instance
	err = newDatabase.Load()
	if err != nil {
		log.Fatalf("Failed to load database: %v", err)
	}

	// Retrieve the data from the new database instance
	username, _ := newDatabase.Get("username")
	email, _ := newDatabase.Get("email")


	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Email: %s\n", email)


	if username != "john_doe" || email != "john@example.com" {
		log.Fatalf("Data did not persist as expected")
	} else {
		fmt.Println("Data persisted successfully.")
	}
}