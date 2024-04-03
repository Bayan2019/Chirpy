package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var ErrNotExist = errors.New("resource does not exist")

// 5. Storage / 1. Storage
// Keep your entire "database" in a single file called database.json
// To make sure
// that multiple requests don't try to write to the database at the same time,
// you should use a mutex to lock the database while you're using it.
type DB struct {
	path string
	mu   *sync.RWMutex
}

// 5. Storage / 1. Storage
// Any time you need to update the database,
// you should read the entire thing into memory (unmarshal it into a struct),
// update the data, and
// then write the entire thing back to disk (marshal it back into JSON).
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	// 5. Storage / 7. Users
	Users map[int]User `json:"users"`
	// Revocations map[string]Revocation `json:"revocations"`
}

// 5. Storage / 1. Storage
// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
		// Revocations: map[string]Revocation{},
	}
	return db.writeDB(dbStructure)
}

// 5. Storage / 1. Storage
// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

// 5. Storage / 1. Storage
// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) ResetDB() error {
	err := os.Remove(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return db.ensureDB()
}

// 5. Storage / 1. Storage
// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure := DBStructure{}

	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}

	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}
