package db

import (
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"os"
	"fmt"
	"strings"
)


// Function for saving db data to file for persistant saves
func (db *db) Save() error {
	data, err := proto.Marshal(db.Database)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(db.filename, data, 0644)
}



func (db *db) Load() error {
	data, err := ioutil.ReadFile(db.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	// Decoding from proocol buffer
	err = proto.Unmarshal(data, db.Database)

	if err != nil {
		return err
	}
	return nil
}


func (sdb *ShardedDB) Save() error {
	sdb.lock.Lock()
	defer sdb.lock.Unlock()
	for _, shard := range sdb.Shards {
		if err := shard.Database.Save(); err != nil {
			return err
		}
	}
	return nil
}


func (sdb *ShardedDB) Load() error {
	sdb.lock.Lock()
	defer sdb.lock.Unlock()
	for _, shard := range sdb.Shards {
		if err := shard.Database.Load(); err != nil {
			return err
		}
	}
	return nil
}





func (wal *WAL) save() {
    data := []byte{}
    for _, entry := range wal.entries {
        data = append(data, []byte(entry+"\n")...)
    }
    if err := ioutil.WriteFile(wal.filename, data, 0644); err != nil {
        fmt.Println("Error saving WAL:", err)
    }
}

func (wal *WAL) load() {
    data, err := ioutil.ReadFile(wal.filename)
    if err != nil {
        if os.IsNotExist(err) {
            return
        }
        fmt.Println("Error loading WAL:", err)
        return
    }
    wal.entries = []string{}
    for _, line := range strings.Split(string(data), "\n") {
        if line != "" {
            wal.entries = append(wal.entries, line)
        }
    }
}