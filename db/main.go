package main

import (
	"fmt"
	"log"
    "bryan/GoDB/dbFiles"
)


func main() {
	shardRanges := [][2]int{
		{0, 9},
		{10, 19},
		{20, 29},
	}

	shardedDB := db.NewShardedDB(shardRanges, "shard_db")

	// Load the sharded database from disk
	if err := shardedDB.Load(); err != nil {
		log.Fatalf("Failed to load database: %v", err)
	}

	// Set some values
	keys := []int{5, 15, 25}
	values := []string{"low", "medium", "high"}

	for i, key := range keys {
		if err := shardedDB.Set(key, values[i]); err != nil {
			log.Fatalf("Failed to set value: %v", err)
		}
	}

	// Get some values
	for _, key := range keys {
		value, err := shardedDB.Get(key)
		if err != nil {
			log.Fatalf("Failed to get value: %v", err)
		}
		fmt.Printf("%d: %s\n", key, value)
	}

	// Save the sharded database to disk
	if err := shardedDB.Save(); err != nil {
		log.Fatalf("Failed to save database: %v", err)
	}
}