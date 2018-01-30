package store

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/bsm/go-guid"
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
		_, err := tx.CreateBucketIfNotExists([]byte(bName))
		return err
	})
	if err != nil {
		return nil, err
	}

	return &BoltDB{FilePath: fp, db: database, bucketName: bName}, nil
}

func (b BoltDB) Insert(m ImageMeta) (string, error) {
	// Define image object to store
	img := image{
		Timestamp: time.Now().Format(time.RFC3339),
		Meta:      m,
	}

	// Encode image
	encoded, err := json.Marshal(img)
	if err != nil {
		return "", err
	}

	// Generate GUID
	id := guid.New96()
	guid := hex.EncodeToString(id.Bytes())

	// Insert into the db
	err = b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucketName))
		return b.Put([]byte(guid), encoded)
	})
	if err != nil {
		return "", err
	}

	return guid, nil
}

func (b BoltDB) Get(guid string) (ImageMeta, error) {
	var img image

	err := b.db.View(func(tx *bolt.Tx) error {
		// Get value
		b := tx.Bucket([]byte(b.bucketName))
		v := b.Get([]byte(guid))

		// Decode JSON
		return json.Unmarshal(v, &img)
	})

	return img.Meta, err
}

func (b BoltDB) Delete(guid string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucketName))
		return b.Delete([]byte(guid))
	})
}

func (b BoltDB) Purge() error {
	// Timestamp an hour ago
	maxTime := time.Now().Add(-time.Hour)

	b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucketName))

		b.ForEach(func(k, v []byte) error {
			var img image
			json.Unmarshal(v, &img)

			imgTime, err := time.Parse(time.RFC3339, img.Timestamp)
			if err != nil {
				return err
			}

			if imgTime.Before(maxTime) {
				return b.Delete(k)
			}

			return nil
		})
		return nil
	})

	return nil
}

func (b BoltDB) Close() error {
	return b.db.Close()
}
