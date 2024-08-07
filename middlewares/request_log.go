package middlewares

import (
	"log"
	"net/http"
	"time"
)

func LogRequests(handler http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
        logger.Printf("new request -> Path: %s - RemoteAddr: %s - Host: %s", r.URL.Path, r.RemoteAddr, r.Host)
		handler.ServeHTTP(w, r)
        logger.Printf("Processed in %d ms", time.Since(start).Milliseconds())
	})
}
