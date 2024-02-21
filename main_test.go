package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/thansetan/pendekin/dto"
	"github.com/thansetan/pendekin/handler"
	"github.com/thansetan/pendekin/helper"
	"github.com/thansetan/pendekin/repository"
	"github.com/thansetan/pendekin/storage"
	"github.com/thansetan/pendekin/usecase"
	"github.com/thansetan/pendekin/util"
)

func TestApp(t *testing.T) {
	dbFileName := "storageTest"
	t.Cleanup(func() {
		runtime.GC()
		os.Remove(dbFileName)
	})
	logger := helper.NewLogger("text")
	urlDB, _ := storage.NewURLDatabase(dbFileName, time.Minute)
	repo := repository.NewURLRepository(urlDB)
	uc := usecase.NewURLUsecase(repo, logger)
	handler := handler.NewURLHandler(uc)

	r := http.NewServeMux()
	r.HandleFunc("GET /", handler.Home(newShortURLHTML))
	r.HandleFunc("GET /{shortURL}", http.HandlerFunc(handler.Get))
	r.HandleFunc("POST /shorten", handler.Save)
	testServer := httptest.NewServer(r)
	defer testServer.Close()
	var (
		originalURL = "https://github.com/thansetan"
		shortURL    string
	)
	t.Run("Home", func(t *testing.T) {
		resp, err := testServer.Client().Get(testServer.URL + "/")
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code to be %d, got %d\n", http.StatusOK, resp.StatusCode)
		}
		if contentType := resp.Header.Get("Content-Type"); contentType != "text/html" {
			t.Errorf("expected content-type to be text/html, got %s\n", contentType)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if string(body) != string(newShortURLHTML) {
			t.Errorf("expected body to be %s, got %s\n", string(newShortURLHTML), string(body))
		}
	})

	t.Run("StoreURL", func(t *testing.T) {
		var buf bytes.Buffer
		reqData := dto.URLRequest{
			OriginalURL: originalURL,
		}
		err := json.NewEncoder(&buf).Encode(reqData)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := testServer.Client().Post(testServer.URL+"/shorten", "application/json", &buf)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected response status code to be %d, got %d\n", http.StatusCreated, resp.StatusCode)
		}
		if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
			t.Errorf("expected content-type to be application/json, got %s\n", contentType)
		}
		defer resp.Body.Close()
		var respData helper.Response[dto.URLResponse]
		err = json.NewDecoder(resp.Body).Decode(&respData)
		if err != nil {
			t.Fatal(err)
		}
		if respData.Data.OriginalURL != reqData.OriginalURL {
			t.Errorf("expected original URL to be %s, got %s\n", reqData.OriginalURL, respData.Data.OriginalURL)
		}
		if len(respData.Data.ShortURL) == 0 {
			t.Fatal("expected short URL to not be empty\n")
		}
		shortURL, _ = util.Shorten(reqData.OriginalURL, len(respData.Data.ShortURL))
		if respData.Data.ShortURL != shortURL {
			t.Fatalf("expeced short URL to be %s, got %s\n", shortURL, respData.Data.ShortURL)
		}
	})

	t.Run("GetURL", func(t *testing.T) {
		client := testServer.Client()
		var finalURL string
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			finalURL = req.URL.String()
			return http.ErrUseLastResponse
		}
		resp, err := client.Get(testServer.URL + "/" + shortURL)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusTemporaryRedirect {
			t.Errorf("expected response status code to be %d, got %d\n", http.StatusTemporaryRedirect, resp.StatusCode)
		}
		if finalURL != originalURL {
			t.Errorf("expected request to redirected to %s, got %s\n", originalURL, finalURL)
		}
	})
}
