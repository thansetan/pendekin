package main

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
		logger.Error("error creating database", "err", err.Error())
		os.Exit(1)
	}
	defer urlDB.Backup()

	repo := repository.NewURLRepository(urlDB)
	uc := usecase.NewURLUsecase(repo, logger)
	handler := handler.NewURLHandler(uc)

	// users/each IP can only create 10 shortlinks/day, will reset at 00.00 UTC every day
	rl := middlewares.NewRateLimiter(10, 24*time.Hour)

	r := http.NewServeMux()
	r.HandleFunc("GET /", handler.Home(newShortURLHTML))
	r.HandleFunc("GET /{slug}", middlewares.GetClientIP(handler.Get))
	r.HandleFunc("POST /shorten", middlewares.GetClientIP(rl.RateLimitMiddleware(handler.Save)))

	srv := new(http.Server)
	srv.Addr = ":8080"
	srv.Handler = r

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("error starting server", "err", err.Error())
			os.Exit(1)
		}
	}()
	fmt.Println("running at http://localhost:8080")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Error("error shutting down server", "err", err.Error())
		os.Exit(1)
	}

	fmt.Println("server is shutting down...")
}
