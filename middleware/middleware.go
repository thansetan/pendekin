package middleware

import (
	"context"
	"net"
	"net/http"

	"github.com/thansetan/pendekin/helper"
)

func GetClientIP(f http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		ctx := context.WithValue(r.Context(), helper.UserIPKey, ip)
		r = r.WithContext(ctx)

		f.ServeHTTP(w, r)
	}
}
