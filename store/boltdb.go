package store

type BoltDB struct {
	FilePath string
}

func (b BoltDB) Init() error {

	return nil
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
