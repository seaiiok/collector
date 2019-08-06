package cache

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/boltdb/bolt"
)

var db *bolt.DB
var once sync.Once

type ICache interface {
	Get(string) string
	Set(string, interface{})
	GetMap() map[string]string
	Close()
}
type cache struct {
	path   string
	bucket string
}

func New(path, bucket string) ICache {
	this := &cache{
		path:   path,
		bucket: bucket,
	}

	once.Do(func() {
		db = &bolt.DB{}
	})

	if db.Stats().FreeAlloc == 0 {
		db, _ = bolt.Open(path, 0600, nil)
	}

	// Start a writable transaction.
	tx, err := db.Begin(true)
	if err != nil {

	}
	defer tx.Rollback()

	// Use the transaction...
	_, err = tx.CreateBucketIfNotExists([]byte(bucket))
	if err != nil {

	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {

	}
	return this
}

func (this *cache) Get(key string) string {
	tx, err := db.Begin(true)
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte(this.bucket))
	if err != nil {
		fmt.Println(err)
	}

	vb := b.Get([]byte(key))

	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	}

	if len(vb) == 0 {
		return ""
	}

	var tempValue string
	err = json.Unmarshal(vb, &tempValue)
	return tempValue
}

func (this *cache) Set(key string, value interface{}) {
	tx, err := db.Begin(true)
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte(this.bucket))
	if err != nil {
		fmt.Println(err)
	}

	vb, _ := json.Marshal(value)
	b.Put([]byte(key), vb)

	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	}
}

func (this *cache) GetMap() (list map[string]string) {
	list = make(map[string]string)
	db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(this.bucket))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var tempValue string
			json.Unmarshal(v, &tempValue)
			list[string(k)] = tempValue
		}

		return nil
	})
	return
}

func (this *cache) Close() {
	if db != nil {
		db.Close()
	}
}
