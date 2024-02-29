package database

import (
	"encoding/json"
	"os"
	"sort"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}

		_, err = file.Write([]byte("{}"))

		if err != nil {
			return nil, err
		}

		err = file.Close()
		if err != nil {
			return nil, err
		}
	}

	return &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	database, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := database.Chirps

	var chirpsSlice []Chirp
	for _, chirp := range chirps {
		chirpsSlice = append(chirpsSlice, chirp)
	}

	sort.Slice(chirpsSlice, func(i, j int) bool {
		return chirpsSlice[i].ID < chirpsSlice[j].ID
	})

	return chirpsSlice, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	chirps, err := db.GetChirps()
	if err != nil {
		return Chirp{}, err
	}

	nextID := len(chirps) + 1
	chirp := Chirp{
		ID:   nextID,
		Body: body,
	}

	chirps[nextID] = chirp

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	dbStructure.Chirps[nextID] = chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	file, err := os.Open(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	decoder := json.NewDecoder(file)
	var dbStructure DBStructure
	err = decoder.Decode(&dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	err = file.Close()
	if err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	file, err := os.Create(db.path)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(dbStructure)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

// // ensureDB creates a new database file if it doesn't exist
// func (db *DB) ensureDB() error {
// 	_, err := os.Stat(db.path)
// 	if os.IsNotExist(err) {
// 		file, err := os.Create(db.path)
// 		if err != nil {
// 			return err
// 		}

// 		_, err = file.Write([]byte("{}"))

// 		if err != nil {
// 			return err
// 		}

// 		err = file.Close()
// 		if err != nil {
// 			return err
// 		}

// 		return nil

// 	}
// 	return nil
// }
