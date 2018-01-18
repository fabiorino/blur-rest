package store

import "fmt"

type ImageMeta struct {
	Blur float32 `json:"user"`
}

func (i *ImageMeta) Validate() error {
	if i.Blur < 0 {
		return fmt.Errorf("Blur must be > 0")
	}

	return nil
}

func (i *ImageMeta) ApplyDefaults() {
	if i.Blur == 0 {
		i.Blur = 0.5
	}
}

type Image struct {
	GUID      string    `json:"guid"`
	Timestamp string    `json:"timestamp"`
	Meta      ImageMeta `json:"meta"`
}

type Store interface {
	Insert(ImageMeta) error
	Get(string) (ImageMeta, error)
	Delete(string) error
	Purge() error
	Close() error
}
