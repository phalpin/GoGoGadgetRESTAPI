package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func SimpleLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("Request Beginning: %v %v", strings.ToUpper(r.Method), r.RequestURI))
		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Now()
		diff := end.Sub(start)
		log.Println(fmt.Sprintf("Request Completed. Time Taken: %v ms, %v µs", diff.Milliseconds(), diff.Microseconds()%1000))
	})
}
