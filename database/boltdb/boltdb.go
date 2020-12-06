package boltdb

import (
	"errors"
	"time"

	"github.com/sjsafranek/node-server/models"
	bolt "go.etcd.io/bbolt"
)

func New(connectionString string) (*BoltDb, error) {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(connectionString, 0600, &bolt.Options{Timeout: 1 * time.Second})
	return &BoltDb{db: db}, err
}

type BoltDb struct {
	db *bolt.DB
}

func (self *BoltDb) Has(bucketName string) bool {
	found := false
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		found = nil != b
		return nil
	})
	return found
}

func (self *BoltDb) Create(bucketName string) (models.Bucket, error) {
	bucket := BoltBucket{bucketName: bucketName, db: self.db}
	return &bucket, self.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(bucketName))
		return err
	})
}

func (self *BoltDb) Delete(bucketName string) error {
	return self.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucketName))
	})
}

func (self *BoltDb) Get(bucketName string) (models.Bucket, error) {
	if !self.Has(bucketName) {
		return nil, errors.New("Not found")
	}
	return &BoltBucket{bucketName: bucketName, db: self.db}, nil
}

type BoltBucket struct {
	bucketName string
	db         *bolt.DB
}

func (self *BoltBucket) view(clbk func(*bolt.Bucket) error) error {
	return self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(self.bucketName))
		return clbk(b)
	})
}

func (self *BoltBucket) update(clbk func(*bolt.Bucket) error) error {
	return self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(self.bucketName))
		return clbk(b)
	})
}

func (self *BoltBucket) Get(key string) ([]byte, error) {
	var value []byte
	return value, self.view(func(b *bolt.Bucket) error {
		v := b.Get([]byte(key))
		value = v
		return nil
	})
}

func (self *BoltBucket) Set(key string, value []byte) error {
	return self.update(func(b *bolt.Bucket) error {
		return b.Put([]byte(key), value)
	})
}

func (self *BoltBucket) Delete(key string) error {
	return self.update(func(b *bolt.Bucket) error {
		return b.Delete([]byte(key))
	})
}
