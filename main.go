package main

import (
	"fmt"
	"go_db/db"
)



func main() {
	// Create a new DB instance
	db := db.NewDb("test.db")

	// Set some key-value pairs
	db.Set("key1", "value1")
	db.Set("key2", "value2")
	db.Set("key3", "value3")
	db.Set("apple", "1")
	db.Set("Bryan", "cool")
	db.Set("bruh", "value3")
	// Save the DB to a file
	err := db.Save()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Load the DB from the file
	err = db.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print out the loaded DB
	for key, value := range db.Store {
		fmt.Printf("%s: %s\n", key, value)
	}
	
}