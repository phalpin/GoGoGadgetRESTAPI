package middleware

import (
	"net/http"
	"time"
)

type RequestMetrics struct {
	CorrelationId       string
	Timestamp           time.Time
	Duration            time.Duration
	RequestMethod       string
	RequestURI          string
	RequestUserAgent    string
	ResponseStatusCode  int
	RequestHeaderValues map[string][]string
}

func MetricsCollector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		respWrite := wrapResponseWriter(w)
		next.ServeHTTP(respWrite, r)
		duration := time.Now().Sub(start)

		metrics := &RequestMetrics{
			CorrelationId:       GetContextStorageValue(r, CorrelationKey),
			Timestamp:           start,
			Duration:            duration,
			RequestMethod:       r.Method,
			RequestURI:          r.RequestURI,
			RequestUserAgent:    r.UserAgent(),
			ResponseStatusCode:  respWrite.status,
			RequestHeaderValues: make(map[string][]string),
		}

		for k, v := range r.Header {
			for _, val := range v {
				metrics.RequestHeaderValues[k] = append(metrics.RequestHeaderValues[k], val)
			}
		}
	})
}
