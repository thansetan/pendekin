package storage_test

import (
	"testing"
	"time"

	data "github.com/thansetan/pendekin/storage"
)

var (
	urlDB    = data.NewURLDatabase("aaa", 5*time.Second)
	longURL1 = "https://google.com"
	longURL2 = "https://x.com"
)

func TestStoreURL(t *testing.T) {

	urlDB.Store("b", longURL2)
	t.Run("test storing a URL", func(t *testing.T) {
		err := urlDB.Store("a", longURL1)
		if err != nil {
			t.Errorf("got an error when storing a URL: %s", err)
		}

	})

}

func TestGetURL(t *testing.T) {
	t.Run("get a short URL", func(t *testing.T) {
		URL, err := urlDB.Get("a")
		if err != nil {
			t.Errorf("got an error: %s", err)
		}

		if URL.OriginalURL != longURL1 {
			t.Errorf("original URL should be %s, got %s instead", longURL1, URL.OriginalURL)
		}
	})

	t.Run("get a invalid URL", func(t *testing.T) {
		_, err := urlDB.Get("xx")
		if err == nil {
			t.Errorf("should return an error, nil instead")
		}

	})

	time.Sleep(4 * time.Second)
	t.Run("get a valid URL after certain time", func(t *testing.T) {
		URL, err := urlDB.Get("b")
		if err != nil {
			t.Errorf("got an error: %s", err)
		}

		if URL.OriginalURL != longURL2 {
			t.Errorf("original URL should be %s, got %s instead", longURL2, URL.OriginalURL)
		}
	})

	time.Sleep(2 * time.Second)
	t.Run("get a deleted URL", func(t *testing.T) {
		_, err := urlDB.Get("a")
		if err == nil {
			t.Errorf("should return an error, nil instead")
		}
	})

	t.Run("get a valid URL after certain time again", func(t *testing.T) {
		URL, err := urlDB.Get("b")
		if err != nil {
			t.Errorf("got an error: %s", err)
		}

		if URL.OriginalURL != longURL2 {
			t.Errorf("original URL should be %s, got %s instead", longURL2, URL.OriginalURL)
		}
	})
}
