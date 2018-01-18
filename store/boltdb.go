package store

import (
	"github.com/boltdb/bolt"
)

type BoltDB struct {
	FilePath   string
	db         *bolt.DB
	bucketName string
}

func NewBoltDB(fp string) (*BoltDB, error) {
	// New database
	database, err := bolt.Open(fp, 0600, nil)
	if err != nil {
		return nil, err
	}

	// New bucket
	bName := "images"
	err = database.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(bName))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &BoltDB{FilePath: fp, db: database, bucketName: bName}, nil
}

func (b BoltDB) Insert(m ImageMeta) error {

	return nil
}

func (b BoltDB) Get(guid string) (ImageMeta, error) {

	return ImageMeta{}, nil
}

func (b BoltDB) Delete(guid string) error {

	return nil
}

func (b BoltDB) Purge() error {

	return nil
}

func (b BoltDB) Close() error {

	return nil
}
