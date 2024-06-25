package db 


import (
	"sync"
	"fmt"
)

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


