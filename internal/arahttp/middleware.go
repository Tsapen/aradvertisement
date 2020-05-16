package arahttp

import (
	"net/http"

	"github.com/Tsapen/aradvertisement/internal/jwt"
)

const (
	origin         = "Access-Control-Allow-Origin"
	allowedOrigin  = "*"
	methods        = "Access-Control-Allow-Methods"
	allowedMethods = "GET, POST, OPTIONS"
	headers        = "Access-Control-Allow-Headers"
	allowedHeaders = "*"
)

func withHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(origin, allowedOrigin)
			w.Header().Set(methods, allowedMethods)
			w.Header().Set(headers, allowedHeaders)
			if r.Method == http.MethodOptions {
				return
			}

			h.ServeHTTP(w, r)
		},
	)
}

func withTokenMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var token = extractToken(r)
			if err := jwt.TokenValid(token); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		},
	)
}
