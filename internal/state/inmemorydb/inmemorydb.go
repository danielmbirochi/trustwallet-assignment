// Package inmemorydb implements the key-value db.
package inmemorydb

import (
	"errors"
	"sync"
)

var (
	// ErrInMemoryDBNotFound is returned if a key is requested that is not found in
	// the provided memory database.
	ErrInMemoryDBNotFound = errors.New("not found")
)

// Database is an ephemeral key-value store.
type Database struct {
	db   map[string][][]byte
	lock sync.RWMutex
}

// New returns a wrapped map with all the required database interface methods
// implemented.
func New() *Database {
	return &Database{
		db: make(map[string][][]byte),
	}
}

// Has retrieves if a key is present in the key-value store.
func (db *Database) Has(key string) (bool, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		return false, ErrInMemoryDBNotFound
	}
	_, ok := db.db[key]
	return ok, nil
}

// Get retrieves the given key if it's present in the key-value store.
func (db *Database) Get(key string) ([][]byte, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		return nil, ErrInMemoryDBNotFound
	}
	if entry, ok := db.db[key]; ok {
		result := make([][]byte, len(entry))
		copy(result, entry)
		return result, nil
	}
	return nil, ErrInMemoryDBNotFound
}

// Put inserts the given value into the key-value store.
func (db *Database) Put(key string, value [][]byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.db == nil {
		return ErrInMemoryDBNotFound
	}
	db.db[key] = append(db.db[key], value...)
	return nil
}

// Delete removes the key from the key-value store.
func (db *Database) Delete(key string) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.db == nil {
		return ErrInMemoryDBNotFound
	}
	delete(db.db, key)
	return nil
}

// Close deallocates the internal map and ensures any consecutive data access op
// fails with an error.
func (db *Database) Close() {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.db = nil
}
