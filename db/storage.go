package db

import (
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"os"
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
