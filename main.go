package main

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/thansetan/pendekin/handler"
	"github.com/thansetan/pendekin/helpers"
	"github.com/thansetan/pendekin/middlewares"
	"github.com/thansetan/pendekin/repository"
	"github.com/thansetan/pendekin/storage"
	"github.com/thansetan/pendekin/usecase"
)

func main() {

	logger := helpers.NewLogger("text")

	// a short link will be deleted automatically 7 days after it was last accessed
	urlDB := storage.NewURLDatabase("urlData", 7*24*time.Hour)
	repo := repository.NewURLRepository(urlDB)
	uc := usecase.NewURLUsecase(repo, logger)
	handler := handler.NewURLHandler(uc)

	// users/each IP can only create 10 shortlinks/day, will reset at 00.00 UTC every day
	rl := middlewares.NewRateLimiter(10, 24*time.Hour)

	http.HandleFunc("/", middlewares.GetClientIP(func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile(`^/([A-Za-z0-9-_]{5})$`)
		if re.MatchString(r.URL.Path) {
			handler.Get(w, r, []byte(r.URL.Path[1:]))
		} else {
			http.NotFound(w, r)
		}
	}))

	http.HandleFunc("/shorten", middlewares.GetClientIP(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			rl.RateLimitMiddleware(handler.Save, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	fmt.Println("running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
