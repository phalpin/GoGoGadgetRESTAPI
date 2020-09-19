package middleware

import (
	"github.com/google/uuid"
	"net/http"
)

const CorrelationKey string = "x-correlation-id"

func CorrelationManager(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		correlationId := r.Header.Get(CorrelationKey)
		if correlationId == "" {
			correlationId = uuid.New().String()
		}

		SetContextStorageValue(r, CorrelationKey, correlationId)
		next.ServeHTTP(w, r)
		w.Header().Add(CorrelationKey, correlationId)
	})
}
