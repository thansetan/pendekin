package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"time"

	"github.com/thansetan/pendekin/handler"
	"github.com/thansetan/pendekin/helpers"
	"github.com/thansetan/pendekin/middlewares"
	"github.com/thansetan/pendekin/repository"
	"github.com/thansetan/pendekin/storage"
	"github.com/thansetan/pendekin/usecase"
)

//go:embed templates/new.html
var newShortURLHTML []byte

func main() {
	logger := helpers.NewLogger("text")

	// a short link will be deleted automatically 7 days after it was last accessed
	urlDB, err := storage.NewURLDatabase("urlData", 7*24*time.Hour)
	if err != nil {
		panic(err)
	}

	repo := repository.NewURLRepository(urlDB)
	uc := usecase.NewURLUsecase(repo, logger)
	handler := handler.NewURLHandler(uc)

	// users/each IP can only create 10 shortlinks/day, will reset at 00.00 UTC every day
	rl := middlewares.NewRateLimiter(10, 24*time.Hour)

    http.HandleFunc("GET /", handler.Home(newShortURLHTML))
	http.HandleFunc("GET /{slug}", middlewares.GetClientIP(handler.Get))
	http.HandleFunc("POST /shorten", middlewares.GetClientIP(rl.RateLimitMiddleware(handler.Save)))

	fmt.Println("running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
