package db
import (
	"sync"
)

type db struct {
	Store map[string]string
	lock sync.RWMutex
	filename string
}

func NewDb(filename string) *db {
	return &db{
		Store: make(map[string]string),
		filename: filename,
	}
}

func (db *db) Set(key, value string) {
	db.lock.Lock()
	defer db.lock.Unlock()
	db.Store[key] = value 
	
}

func (db *db) Get(key string) (string, bool) {
	db.lock.RLock()
	defer db.lock.RUnlock()
	value, exists := db.Store[key]
	return value, exists

}

func (db *db) Delete(key string) {
	db.lock.Lock()
	defer db.lock.Unlock()
	delete(db.Store, key)
}


