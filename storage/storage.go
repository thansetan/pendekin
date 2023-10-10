package storage

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

type URL struct {
	OriginalURL string
	deleteCh    <-chan time.Time
}

type URLData struct {
	data        map[string]URL
	mu          *sync.RWMutex
	deleteAfter time.Duration
	fileName    string
}

func NewURLDatabase(fileName string, deleteAfter time.Duration) *URLData {
	urlData := &URLData{
		mu:          new(sync.RWMutex),
		deleteAfter: deleteAfter,
		fileName:    fileName,
	}

	file, _ := os.Open(fileName)
	if file != nil {
		err := gob.NewDecoder(file).Decode(&urlData.data)
		if err != nil {
			fmt.Println(err)
		}
		go urlData.delete()
		return urlData
	}

	urlData.data = make(map[string]URL)
	go urlData.delete()
	return urlData
}

func (d *URLData) Store(shortURL, longURL string) error {
	d.mu.Lock()
	defer d.backup()
	defer d.mu.Unlock()
	if URL, ok := d.data[shortURL]; ok {
		if URL.OriginalURL != longURL {
			return errors.New("something wrong")
		}
		return nil
	}

	d.data[shortURL] = URL{
		OriginalURL: longURL,
		deleteCh:    time.After(d.deleteAfter),
	}
	return nil
}

func (d *URLData) Get(shortURL string) (URL, error) {
	d.mu.Lock()
	defer d.backup()
	defer d.mu.Unlock()
	if data, ok := d.data[shortURL]; !ok {
		return data, fmt.Errorf("key %s not found", shortURL)
	} else {
		data.deleteCh = time.After(d.deleteAfter)
		d.data[shortURL] = data
		return data, nil
	}
}

func (d *URLData) delete() {
	for {
		d.mu.Lock()
		for key, url := range d.data {
			select {
			case <-url.deleteCh:
				delete(d.data, key)
			default:
				continue
			}
		}
		d.mu.Unlock()
	}
}

func (d *URLData) backup() {
	file, err := os.OpenFile(d.fileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	err = gob.NewEncoder(file).Encode(d.data)
}
