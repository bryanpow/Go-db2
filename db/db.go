package db
import (

	"sync"
)

type db struct {
	store map[string]string
	lock sync.RWMutex
}

func NewDb() *db {
	return &db{
		store: make(map[string]string),
	}
}

func (db *db) Set(key, value string) {
	db.lock.Lock()
	db.store[key] = value 
	db.lock.Unlock()
}

func (db *db) Get(key string) (string, bool) {
	db.lock.RLock()
	value, exists := db.store[key]
	db.lock.RUnlock()
	return value, exists

}

func (db *db) Delete(key string) {
	db.lock.Lock()
	delete(db.store, key)
	db.lock.Unlock()
}
