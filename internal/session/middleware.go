package session

import (
	"context"
	"net/http"
	"slices"
)

func GetToken(req *http.Request) *Token {
	v := req.Context().Value(tokenKey{})
	if v == nil {
		return nil
	}

	return v.(*Token)
}

func Middleware(store *Store, methods ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			token, err := store.Get(req, rw)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			if slices.Contains(methods, req.Method) {
				store.Set(rw, token)
			}

			ctx := context.WithValue(req.Context(), tokenKey{}, token)
			next.ServeHTTP(rw, req.WithContext(ctx))
		})
	}
}

type tokenKey struct{}
