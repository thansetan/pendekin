package middleware

import (
	"log/slog"
	"net/http"
)

func Recover(f http.Handler, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if err == http.ErrAbortHandler {
					panic(err)
				}
				logger.Error("panic recovered", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		f.ServeHTTP(w, r)
	}
}
