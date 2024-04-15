package db
import (
	"sync"
)


//  Struct for database.
//  Lock for concurrent actions
//  Filename for saving to disk (add prototype)
type db struct {
	Store map[string]string
	lock sync.RWMutex
	filename string
}


//Fucntion for making new database
func NewDb(filename string) *db {
	return &db{
		Store: make(map[string]string),
		filename: filename,
	}
}


// Function for setting new item in the database
func (db *db) Set(key, value string) {
	db.lock.Lock()
	defer db.lock.Unlock()
	db.Store[key] = value 
	
}

// Function for getting item from the database
func (db *db) Get(key string) (string, bool) {
	db.lock.RLock()
	defer db.lock.RUnlock()
	value, exists := db.Store[key]
	return value, exists

}


// Function for deleting item in database
func (db *db) Delete(key string) {
	db.lock.Lock()
	defer db.lock.Unlock()
	delete(db.Store, key)
}


