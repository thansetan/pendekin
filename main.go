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

	// unaccessed short link will be deleted after 7 days after its creation
	urlDB := storage.NewURLDatabase("urlData", 7*24*time.Hour)
	repo := repository.NewURLRepository(urlDB)
	uc := usecase.NewURLUsecase(repo, logger)
	handler := handler.NewURLHandler(uc)

	// users can only create 10 shortlinks/day, will reset at 00.00 UTC every day
	rl := middlewares.NewRateLimiter(3, time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC))

	http.HandleFunc("/", middlewares.GetClientIP(func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile(`\/([A-Za-z0-9-_]+)\b`)
		path := re.Find([]byte(r.URL.Path))
		switch len(path) {
		case 0:
			helpers.ResponseBuilder(w, http.StatusOK, "", "hello world")
		case 6:
			handler.Get(w, r, path[1:])
		default:
			helpers.ResponseBuilder(w, http.StatusNotFound, "page not found", nil)
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

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%v", rl.GetUsers())
	})

	fmt.Println("running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
