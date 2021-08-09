package database

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v3"
)

var instance *DB

type DB struct {
	*badger.DB
}

func New() (*DB, error) {
	if instance == nil {
		opts := badger.DefaultOptions(Path()).WithLoggingLevel(badger.WARNING)
		db, err := badger.Open(opts)
		if err != nil {
			return nil, err
		}

		instance = &DB{db}
	}

	return instance, nil
}

func Close() error {
	if instance != nil {
		return instance.Close()
	}

	return fmt.Errorf("nothing to close")
}

func (db *DB) Set(key, value []byte) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (db *DB) Get(key []byte) ([]byte, error) {
	val := make([]byte, 0)
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		val, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return []byte{}, err
	}
	return val, nil
}
