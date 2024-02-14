package storage

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

type url struct {
	OriginalURL  string
	LastAccessed time.Time
}

type URLData struct {
	data        map[string]url
	mu          *sync.RWMutex
	deleteAfter time.Duration
	fileName    string
}

func NewURLDatabase(fileName string, deleteAfter time.Duration) (*URLData, error) {
	urlData := &URLData{
		mu:          new(sync.RWMutex),
		deleteAfter: deleteAfter,
		fileName:    fileName,
	}

	file, _ := os.Open(fileName)
	if file != nil {
		err := gob.NewDecoder(file).Decode(&urlData.data)
		if err != nil {
			return urlData, err
		}
		go urlData.delete()
		return urlData, nil
	}
	urlData.data = make(map[string]url)
	go urlData.delete()
	return urlData, nil
}

func (d *URLData) Store(shortURL, longURL string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if URL, ok := d.data[shortURL]; ok {
		if URL.OriginalURL != longURL {
			return errors.New("something wrong")
		}
		return nil
	}

	d.data[shortURL] = url{
		OriginalURL:  longURL,
		LastAccessed: time.Now().UTC(),
	}
	return nil
}

func (d *URLData) Get(shortURL string) (url, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if data, ok := d.data[shortURL]; !ok {
		return data, fmt.Errorf("key %s not found", shortURL)
	} else {
		data.LastAccessed = time.Now().UTC()
		d.data[shortURL] = data
		return data, nil
	}
}

func (d *URLData) delete() {
	for {
		d.mu.Lock()
		for key, url := range d.data {
			if time.Since(url.LastAccessed) > d.deleteAfter {
				delete(d.data, key)
			}
		}
		d.mu.Unlock()
	}
}

func (d *URLData) Backup() {
	file, err := os.OpenFile(d.fileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	err = gob.NewEncoder(file).Encode(d.data)
	if err != nil {
		fmt.Println(err)
	}
}
