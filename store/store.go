package store

import "fmt"

type ImageMeta struct {
	Blur int `json:"blur"`
}

func (i *ImageMeta) Validate() error {
	if i.Blur < 0 {
		return fmt.Errorf("Blur must be > 0")
	}

	return nil
}

func (i *ImageMeta) ApplyDefaults() {
	if i.Blur == 0 {
		i.Blur = 10
	}
}

type image struct {
	Timestamp string
	Meta      ImageMeta
}

type Store interface {
	Insert(ImageMeta) (string, error)
	Get(string) (ImageMeta, error)
	Delete(string) error
	Purge() error
	Close() error
}
