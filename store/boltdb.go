package store

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/bsm/go-guid"
)

// BoltDB represents the information needed to interact with a BoltDB database
type BoltDB struct {
	FilePath   string
	db         *bolt.DB
	bucketName string
}

// NewBoltDB returns a new BoltDB instance
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

// Insert adds a new image metadata to the db
func (b BoltDB) Insert(m ImageMeta) (string, error) {
	// Define image object to store
	img := image{
		Timestamp: time.Now().Format(time.RFC3339),
		Meta:      m,
	}

	// Encode metadata
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

// Get returns the metadata of an image given its guid, if exists
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

// Delete removes an entry from the db
func (b BoltDB) Delete(guid string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucketName))
		return b.Delete([]byte(guid))
	})
}

// Close closes the BoltDB session
func (b BoltDB) Close() error {
	return b.db.Close()
}
