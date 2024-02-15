package storage

import (
	"os"
	"runtime"
	"testing"
	"time"
)

var (
	testStorageName = "test"
	longURL1        = "https://google.com"
	longURL2        = "https://x.com"
)

func TestNewStorage(t *testing.T) {
	var newURLDB *urlDB
	t.Cleanup(func() {
		newURLDB.SaveToDrive() // backup for next test
	})

	t.Run("create new DB instance (new storage)", func(t *testing.T) {
		var err error
		newURLDB, err = NewURLDatabase(testStorageName, 5*time.Second)
		if err != nil {
			t.Errorf("got an error : %s", err)
		}

		if newURLDB == nil {
			t.Error("expected newURLDB to be *URLData, got nil instead")
		}
	})

	t.Run("store a URL", func(t *testing.T) {
		err := newURLDB.Store("a", longURL1)
		if err != nil {
			t.Errorf("got an error : %s", err)
		}

		err = newURLDB.Store("b", longURL2)
		if err != nil {
			t.Errorf("got an error : %s", err)
		}

	})

	t.Run("get a valid short URL", func(t *testing.T) {
		URL, err := newURLDB.Get("a")
		if err != nil {
			t.Errorf("got an error: %s", err)
		}

		if URL.OriginalURL != longURL1 {
			t.Errorf("original URL should be %s, got %s instead", longURL1, URL.OriginalURL)
		}
	})

	t.Run("get an invalid URL", func(t *testing.T) {
		_, err := newURLDB.Get("xx")
		if err == nil {
			t.Errorf("should return an error, nil instead")
		}

	})

	time.Sleep(4 * time.Second)
	t.Run("get a valid URL after certain time", func(t *testing.T) {
		URL, err := newURLDB.Get("b")
		if err != nil {
			t.Errorf("got an error: %s", err)
		}

		if URL.OriginalURL != longURL2 {
			t.Errorf("original URL should be %s, got %s instead", longURL2, URL.OriginalURL)
		}
	})

	time.Sleep(2 * time.Second)
	t.Run("get a deleted URL", func(t *testing.T) {
		_, err := newURLDB.Get("a")
		if err == nil {
			t.Errorf("should return an error, nil instead")
		}
	})

	t.Run("get a valid URL after certain time again", func(t *testing.T) {
		URL, err := newURLDB.Get("b")
		if err != nil {
			t.Errorf("got an error: %s", err)
		}

		if URL.OriginalURL != longURL2 {
			t.Errorf("original URL should be %s, got %s instead", longURL2, URL.OriginalURL)
		}
	})
}

func TestLoadStorage(t *testing.T) {
	var loadURLDB *urlDB
	t.Cleanup(func() {
		runtime.GC()
		os.Remove(testStorageName)
	})

	t.Run("create new DB instance (load storage)", func(t *testing.T) {
		var err error
		loadURLDB, err = NewURLDatabase(testStorageName, 2*time.Second)
		if err != nil {
			t.Errorf("got an error : %s", err)
		}

		if loadURLDB == nil {
			t.Error("expected loadURLDB to be *URLData, got nil instead")
		}
	})

	t.Run("on a new DB instance, get URL", func(t *testing.T) {
		URL, err := loadURLDB.Get("b")
		if err != nil {
			t.Errorf("got an error: %s", err)
		}

		if URL.OriginalURL != longURL2 {
			t.Errorf("original URL should be %s, got %s instead", longURL2, URL.OriginalURL)
		}
	})

	time.Sleep(3 * time.Second)
	t.Run("on new DB instance, get a (previously) valid URL after URL deleted because of inactivity/not accessed after a period of time", func(t *testing.T) {
		_, err := loadURLDB.Get("b")
		if err == nil {
			t.Errorf("should return an error, nil instead")
		}
	})

	t.Run("store a URL", func(t *testing.T) {
		err := loadURLDB.Store("a", longURL1)
		if err != nil {
			t.Errorf("got an error : %s", err)
		}
	})

}

// var DB *urlDB

// func init() {
// 	DB = &urlDB{
// 		data:        make(map[string]url),
// 		mu:          new(sync.RWMutex),
// 		deleteAfter: 5 * time.Second,
// 		fileName:    "bench",
// 	}
// 	DB.StoreNoBackup("a", longURL1)
// }

// func BenchmarkGet(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		_, _ = DB.Get("a")
// 	}
// }

// func BenchmarkGetNoBackup(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		_, _ = DB.GetNoBackup("a")
// 	}
// }

// func BenchmarkStore(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		_ = DB.Store("b", longURL2)
// 	}
// }

// func BenchmarkStoreNoBackup(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		_ = DB.StoreNoBackup("b", longURL2)
// 	}
// }
