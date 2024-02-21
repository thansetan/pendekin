package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thansetan/pendekin/handler"
	"github.com/thansetan/pendekin/helper"
	"github.com/thansetan/pendekin/middleware"
	"github.com/thansetan/pendekin/repository"
	"github.com/thansetan/pendekin/storage"
	"github.com/thansetan/pendekin/usecase"
)

//go:embed templates/new.html
var newShortURLHTML []byte

func main() {
	logger := helper.NewLogger("text")

	// a short link will be deleted automatically 7 days after it was last accessed
	urlDB, err := storage.NewURLDatabase("urlData", 7*24*time.Hour)
	if err != nil {
		logger.Error("error creating database", "err", err.Error())
		os.Exit(1)
	}
	defer urlDB.SaveToDrive()

	repo := repository.NewURLRepository(urlDB)
	uc := usecase.NewURLUsecase(repo, logger)
	handler := handler.NewURLHandler(uc)

	// users/each IP can only create 10 shortlinks/day
	rl := middleware.NewRateLimiter(10, 24*time.Hour, helper.UserIPKey)

	r := http.NewServeMux()
	r.HandleFunc("GET /", handler.Home(newShortURLHTML))
	r.HandleFunc("GET /{shortURL}", middleware.GetClientIP(http.HandlerFunc(handler.Get)))
	r.HandleFunc("POST /shorten", middleware.GetClientIP(rl.RateLimitMiddleware(handler.Save)))

	srv := new(http.Server)
	srv.Addr = "0.0.0.0:8080"
	srv.Handler = middleware.Recover(r, logger)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("error starting server", "err", err.Error())
			os.Exit(1)
		}
	}()
	fmt.Printf("running at http://%s\n", srv.Addr)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Error("error shutting down server", "err", err.Error())
		os.Exit(1)
	}

	fmt.Println("server is shutting down...")
}
