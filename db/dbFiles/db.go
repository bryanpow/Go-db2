package db

import (
	"errors"
	"sync"
	"fmt"
	"strings"
)

//  Struct for database.
//  Lock for concurrent actions
//  Filename for saving to disk
//  WAL for logging actions
type db struct {
	*Database
	lock sync.RWMutex
	filename string
	wal *WAL
}






//Function for making new database
func NewDb(filename string) *db {
    walFilename := filename + "_wal"
    return &db{
        Database: &Database{
            Store: make(map[string]string),
        },
        filename: filename,
        wal:      NewWAL(walFilename),
    }
}




//function for setting item in shard
func (db *db) Set(key, value string) error {
    db.lock.Lock()
    defer db.lock.Unlock()
    db.wal.Append(fmt.Sprintf("SET %s %s", key, value))
    db.Store[key] = value
    return nil
}


// // Function for getting item from the database
func (db *db) Get(key string) (string, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()
	value, exists := db.Store[key]
	if !exists {
		return "", errors.New("item not found in database")
	}
	return value, nil
}
// Function for deleting item in database

func (db *db) Delete(key string) error {
    db.lock.Lock()
    defer db.lock.Unlock()
    if _, exists := db.Store[key]; !exists {
        return errors.New("item does not exist")
    }
    db.wal.Append(fmt.Sprintf("DELETE %s", key))
    delete(db.Store, key)
    return nil
}


func (db *db) Recover() error {
    db.lock.Lock()
    defer db.lock.Unlock()
    entries := db.wal.GetEntries()
    fmt.Println("Recovering from WAL entries:", entries) // Debug log
    for _, entry := range entries {
        parts := strings.Split(entry, " ")
        if parts[0] == "SET" {
            db.Store[parts[1]] = parts[2]
        } else if parts[0] == "DELETE" {
            delete(db.Store, parts[1])
        }
    }
    return nil
}




