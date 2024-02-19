package repository

import "github.com/thansetan/pendekin/storage"

type URLRepository interface {
	Save(string, string) error
	Get(string) (string, error)
}

type urlRepositoryImpl struct {
	urlData storage.KVStorage
}

func NewURLRepository(urlData storage.KVStorage) *urlRepositoryImpl {
	return &urlRepositoryImpl{
		urlData: urlData,
	}
}

func (repo *urlRepositoryImpl) Save(longURL, shortURL string) error {
	return repo.urlData.Store(shortURL, longURL)
}

func (repo *urlRepositoryImpl) Get(shortURL string) (string, error) {
	var url string

	URL, err := repo.urlData.Get(shortURL)
	if err != nil {
		return url, err
	}

	url = URL.OriginalURL

	return url, nil
}
