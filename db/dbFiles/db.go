package db

import (
	"errors"
	"sync"
	"fmt"
	"strconv"
	"time"
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


type Shard struct {
	ID int
	Range [2]int
	Database *db
	Replicas []*db
}


type ShardedDB struct {
	Shards []*Shard
	lock sync.RWMutex
}

type WAL struct {
    entries  []string
    lock     sync.Mutex
    filename string
}



func NewWAL(filename string) *WAL {
    wal := &WAL{
        entries:  make([]string, 0),
        filename: filename,
    }
    wal.load()
    return wal
}

func (wal *WAL) Append(entry string) {
    wal.lock.Lock()
    defer wal.lock.Unlock()
    fmt.Println("Appending to WAL:", entry) // Debug log
    wal.entries = append(wal.entries, entry)
    wal.save()
}

func (wal *WAL) GetEntries() []string {
    wal.lock.Lock()
    defer wal.lock.Unlock()
    return wal.entries
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


func NewShardedDB(sharedRanges [][2]int, basedFilename string, replicas int) *ShardedDB {
    shardedDb := &ShardedDB{}
    for id, r := range sharedRanges {
        filename := fmt.Sprintf("%s_%d", basedFilename, id)
        shard := &Shard{
            ID: id,
            Range: r,
            Database: NewDb(filename),
        }
        // Add replicas
        for i := 0; i < replicas; i++ {
            replicaFilename := fmt.Sprintf("%s_%d_replica_%d", basedFilename, id, i)
            shard.Replicas = append(shard.Replicas, NewDb(replicaFilename))
        }
        shardedDb.Shards = append(shardedDb.Shards, shard)
    }
    return shardedDb
}







func (sdb *ShardedDB) getShard(key int) (*Shard, error) {
	for _, shard := range sdb.Shards {
		if key >= shard.Range[0] && key <= shard.Range[1] {
			return shard, nil
		}
	}
	fmt.Printf("Key '%d' not found in any shard range\n", key)
	return nil, errors.New("no shard found for key")
}


//function for setting item in shardedDb
func (sdb *ShardedDB) Set(key int, value string) error {
	sdb.lock.Lock()
	defer sdb.lock.Unlock()
	shard, err := sdb.getShard(key)
	if err != nil {
		return err
	}
	if err := shard.Database.Set(strconv.Itoa(key), value); err != nil {
		return err
	}

	// Replicate to all replicas
	for _, replica := range shard.Replicas {
		if err := replica.Set(strconv.Itoa(key), value); err != nil {
			return err
		}
	}
	return nil
}



// Function for getting an item from a shard
func (sdb *ShardedDB) Get(key int) (string, error) {
	sdb.lock.RLock()
	defer sdb.lock.RUnlock()
	shard, err := sdb.getShard(key)
	if err != nil {
		return "", err
	}
	return shard.Database.Get(strconv.Itoa(key))
}


func (sdb *ShardedDB) Delete(key int) error {
	sdb.lock.Lock()
	defer sdb.lock.Unlock()
	shard, err := sdb.getShard(key)
	if err != nil {
		return err
	}
	if err := shard.Database.Delete(strconv.Itoa(key)); err != nil {
		return err
	}
	// Replicate to all replicas
	for _, replica := range shard.Replicas {
		if err := replica.Delete(strconv.Itoa(key)); err != nil {
			return err
		}
	}
	return nil
}



func (sdb *ShardedDB) PromoteReplica(shard *Shard, replica *db) {
	sdb.lock.Lock()
	defer sdb.lock.Unlock()
	shard.Database = replica
	// Remove promoted replica from replicas list
	newReplicas := []*db{}
	for _, rep := range shard.Replicas {
		if rep != replica {
			newReplicas = append(newReplicas, rep)
		}
	}
	shard.Replicas = newReplicas
}


// Function to check if a shard has failed
func shardFailed(shard *db) bool {
    // Example: try to read a known key to check if the shard is responding
    _, err := shard.Get("healthcheck")
    return err != nil
}


// Function to detect and handle shard failures 
func (sdb *ShardedDB) MonitorShards() {
	for _, shard := range sdb.Shards {
		// Simulated check for shard availability
		if shardFailed(shard.Database) {
			if len(shard.Replicas) > 0 {
				sdb.PromoteReplica(shard, shard.Replicas[0])
			} else {
				// Handle no replicas available scenario
				fmt.Println("No replicas available for shard:", shard.ID)
			}
		}
	}
}

func (sdb *ShardedDB) StartMonitoring() {
	ticker := time.NewTicker(30 * time.Second) 
	defer ticker.Stop()
	for range ticker.C {
		sdb.MonitorShards()
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




