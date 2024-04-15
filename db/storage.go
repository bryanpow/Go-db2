package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
)


// Function for saving db data to file for persistant saves
func (db *db) Save() error {
	data, err := json.Marshal(db.Store)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(db.filename, data, 0644)
}


// Function for loading persistent data from json file, and populating the key value store with that data
func (db *db) Load() error {
	data, err := ioutil.ReadFile(db.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	// Decoding to json
	err = json.Unmarshal(data, &db.Store)

	if err != nil {
		return err
	}
	return nil
}
