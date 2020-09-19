package middleware

import (
	"context"
	"net/http"
)

type ContextStorageKey struct{}

func ContextStorage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(r.Context(), &ContextStorageKey{}, make(map[string]string))
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

func getContextStorage(r *http.Request) map[string]string {
	itemMap, _ := r.Context().Value(&ContextStorageKey{}).(map[string]string)
	return itemMap
}

func GetContextStorageValue(r *http.Request, key string) string {
	items := getContextStorage(r)
	return items[key]
}

func SetContextStorageValue(r *http.Request, key string, value string) {
	items := getContextStorage(r)
	items[key] = value
}
