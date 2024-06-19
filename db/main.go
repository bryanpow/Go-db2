package main

import (
	 "encoding/json"
    "log"
    "net/http"
    "strconv"
    "sync"
    "bryan/GoDB/dbFiles"
    "github.com/gorilla/mux"
    "path/filepath"
    "os"
    "os/signal"
    "syscall"
    "context"
    "time"
   
)

var (
    shardedDBInstances = make(map[string]*db.ShardedDB)
    dbMutex            sync.Mutex
)


type Response struct {
	Message string `json:"message"`
}



func loadUserIDs() []string {
    var userIDs []string
    files, err := filepath.Glob("*.db")
    if err != nil {
        log.Fatalf("Failed to load user IDs: %v", err)
    }
    for _, file := range files {
        userID := file[:len(file)-3] // Assuming ".db" extension
        userIDs = append(userIDs, userID)
    }
    return userIDs
}

func getUserShardedDB(userID string) *db.ShardedDB {
    dbMutex.Lock()
    defer dbMutex.Unlock()

    if _, exists := shardedDBInstances[userID]; !exists {
        shardedDBInstances[userID] = db.NewShardedDB([][2]int{{0, 100}, {101, 200}, {201, 300}}, userID+"_db", 2)
        shardedDBInstances[userID].Load() // Load data from files
    }
    return shardedDBInstances[userID]
}

func main() {
    userIDs := loadUserIDs()
    for _, userID := range userIDs {
        getUserShardedDB(userID)
    }

    router := mux.NewRouter()

    router.HandleFunc("/api/{userID}/createdb", createShardedDBHandler).Methods("POST")
    router.HandleFunc("/api/{userID}/set/{key}/{value}", setHandler).Methods("POST")
    router.HandleFunc("/api/{userID}/get/{key}", getHandler).Methods("GET")
    router.HandleFunc("/api/{userID}/delete/{key}", deleteHandler).Methods("DELETE")

    srv := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    // Start monitoring shards in a separate goroutine
    go startMonitoringAllShards()

    go func() {
        log.Println("Server starting on port 8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("ListenAndServe(): %v", err)
        }
    }()

    gracefulShutdown(srv)
}



func startMonitoringAllShards() {
    for {
        dbMutex.Lock()
        for _, shardedDB := range shardedDBInstances {
            go shardedDB.StartMonitoring()
        }
        dbMutex.Unlock()
        time.Sleep(30 * time.Second)
    }
}

func gracefulShutdown(srv *http.Server) {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    <-c

    log.Println("Shutting down gracefully...")

    dbMutex.Lock()
    for _, shardedDB := range shardedDBInstances {
        if err := shardedDB.Save(); err != nil {
            log.Printf("Failed to save database: %v", err)
        }
    }
    dbMutex.Unlock()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    log.Println("Server exiting")
}



func createShardedDBHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["userID"]

    dbMutex.Lock()
    defer dbMutex.Unlock()

    if _, exists := shardedDBInstances[userID]; exists {
        http.Error(w, "Database already exists for user", http.StatusBadRequest)
        return
    }

    shardedDBInstances[userID] = db.NewShardedDB([][2]int{{0, 100}, {101, 200}, {201, 300}}, userID+"_db", 2)
    shardedDBInstances[userID].Save() // Save initial state to disk
    json.NewEncoder(w).Encode(Response{Message: "Sharded database created successfully"})
}

func setHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["userID"]
    keyStr := vars["key"]
    value := vars["value"]

    key, err := strconv.Atoi(keyStr)
    if err != nil {
        http.Error(w, "Invalid key", http.StatusBadRequest)
        return
    }

    userDB := getUserShardedDB(userID)
    err = userDB.Set(key, value)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    userDB.Save() // Save after setting a value
    json.NewEncoder(w).Encode(Response{Message: "Key set successfully"})
}

func getHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["userID"]
    keyStr := vars["key"]

    key, err := strconv.Atoi(keyStr)
    if err != nil {
        http.Error(w, "Invalid key", http.StatusBadRequest)
        return
    }

    userDB := getUserShardedDB(userID)
    value, err := userDB.Get(key)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(Response{Message: value})
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["userID"]
    keyStr := vars["key"]

    key, err := strconv.Atoi(keyStr)
    if err != nil {
        http.Error(w, "Invalid key", http.StatusBadRequest)
        return
    }

    userDB := getUserShardedDB(userID)
    err = userDB.Delete(key)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    userDB.Save() // Save after deleting a value
    json.NewEncoder(w).Encode(Response{Message: "Key deleted successfully"})
}