package db

import (
	"testing"
	"os"
)


func TestSetAndGet(t *testing.T) {
	shardRanges := [][2]int{
		{0, 9},
		{10, 19},
		{20, 29},
	}

	
	shardedDB := NewShardedDB(shardRanges, "shard_db", 2)

	// Set values
	testCases := []struct {
		key   int
		value string
	}{
		{5, "low"},
		{15, "medium"},
		{25, "high"},
	}

	for _, tc := range testCases {
		err := shardedDB.Set(tc.key, tc.value)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}
	}

	// Get values and verify
	for _, tc := range testCases {
		value, err := shardedDB.Get(tc.key)
		if err != nil {
			t.Fatalf("Failed to get value for key %d: %v", tc.key, err)
		}
		if value != tc.value {
			t.Fatalf("Expected value %s for key %d, got %s", tc.value, tc.key, value)
		}
	}
}


func TestDelete(t *testing.T) {
	shardRanges := [][2]int{
		{0, 9},
		{10, 19},
		{20, 29},
	}

	shardedDB := NewShardedDB(shardRanges, "shard_db", 2)

	// Set a value
	err := shardedDB.Set(15, "medium")
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	// Delete the value
	err = shardedDB.Delete(15)
	if err != nil {
		t.Fatalf("Failed to delete value: %v", err)
	}

	// Verify deletion
	_, err = shardedDB.Get(15)
	if err == nil {
		t.Fatalf("Expected error for getting deleted key, got nil")
	}
	if err.Error() != "item not found in database" {
		t.Fatalf("Expected 'item not found in database' error, got %v", err)
	}
}


// Testing setting and getting values in different shards
func TestSharding(t *testing.T) {
    shardRanges := [][2]int{
        {0, 9},
        {10, 19},
        {20, 29},
    }

    shardedDB := NewShardedDB(shardRanges, "shard_db", 2)

    // Set values in different shards
    testCases := []struct {
        key   int
        value string
    }{
        {5, "low"},
        {15, "medium"},
        {25, "high"},
    }

    for _, tc := range testCases {
        err := shardedDB.Set(tc.key, tc.value)
        if err != nil {
            t.Fatalf("Failed to set value: %v", err)
        }
    }

    // Get values and verify they are in the correct shards
    for _, tc := range testCases {
        value, err := shardedDB.Get(tc.key)
        if err != nil {
            t.Fatalf("Failed to get value for key %d: %v", tc.key, err)
        }
        if value != tc.value {
            t.Fatalf("Expected value %s for key %d, got %s", tc.value, tc.key, value)
        }
    }
}



func TestWriteAheadLog(t *testing.T) {
    // Clean up any existing test files
    os.Remove("test_database_file")
    os.Remove("test_database_file_wal")

    // Initialize the database
    db := NewDb("test_database_file")

    // Perform some operations
    if err := db.Set("key1", "value1"); err != nil {
        t.Fatalf("Failed to set key1: %v", err)
    }
    if err := db.Set("key2", "value2"); err != nil {
        t.Fatalf("Failed to set key2: %v", err)
    }
    if err := db.Delete("key1"); err != nil {
        t.Fatalf("Failed to delete key1: %v", err)
    }

    // Simulate a crash and recovery
    recoveredDb := NewDb("test_database_file")
    if err := recoveredDb.Recover(); err != nil {
        t.Fatalf("Failed to recover database: %v", err)
    }

    // Check recovered state
    value, err := recoveredDb.Get("key2")
    if err != nil {
        t.Fatalf("Failed to get key2: %v", err)
    }
    if value != "value2" {
        t.Errorf("Expected value2, got %s", value)
    }

    // Check that key1 is deleted
    _, err = recoveredDb.Get("key1")
    if err == nil {
        t.Errorf("Expected error for getting deleted key1, got nil")
    }

    // Clean up test files
    os.Remove("test_database_file")
    os.Remove("test_database_file_wal")
}