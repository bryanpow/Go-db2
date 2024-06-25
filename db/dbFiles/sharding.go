package db
import (
	"errors"
	"sync"
	"fmt"
	"strconv"
	"time"
)


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

