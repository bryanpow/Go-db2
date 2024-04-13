package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func (db *db) Save() error {
	data, err := json.Marshal(db.Store)
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

	err = json.Unmarshal(data, &db.Store)

	if err != nil {
		return err
	}
	return nil
}