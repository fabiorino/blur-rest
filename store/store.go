package store

import "fmt"

// ImageMeta represents the parameters the user can provide in their POST request
type ImageMeta struct {
	Blur int `json:"blur"`
}

// Validate checks the parameters the user provided
func (i *ImageMeta) Validate() error {
	if i.Blur < 0 {
		return fmt.Errorf("Blur must be > 0")
	}

	return nil
}

// ApplyDefaults sets a default value for the unset parameters
func (i *ImageMeta) ApplyDefaults() {
	if i.Blur == 0 {
		i.Blur = 10
	}
}

// The structure to store in the database
type image struct {
	Timestamp string    `json:"timestamp"`
	Meta      ImageMeta `json:"meta"`
}

// Store represents the functions that every db must implement to work with this API
type Store interface {
	Insert(ImageMeta) (string, error)
	Get(string) (ImageMeta, error)
	Delete(string) error
	Close() error
}
