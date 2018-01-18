package store

import (
	"encoding/hex"
	"encoding/json"

	"github.com/boltdb/bolt"
	guid "github.com/bsm/go-guid"
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
	// Encode image meta
	encoded, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// Generate GUID
	id := guid.New96()
	guid := hex.EncodeToString(id.Bytes())

	// Insert into the db
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucketName))
		err := b.Put([]byte(guid), encoded)
		return err
	})
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
