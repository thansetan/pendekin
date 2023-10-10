package middlewares

import (
	"context"
	"net"
	"net/http"
)

func GetClientIP(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		ctx := context.WithValue(r.Context(), "user_ip", ip)
		r = r.WithContext(ctx)

		f.ServeHTTP(w, r)
	}
}
