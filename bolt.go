package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func openDB() *bolt.DB {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func initBucket(db *bolt.DB, name string) {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func writeToDB(db *bolt.DB, key, value, bucket string) {
	initBucket(db, bucket)
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
}

func readFromDB(db *bolt.DB, key, bucket string) []byte {
	v := []byte("")
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v = b.Get([]byte(key))
		return nil
	})

	return v

}
