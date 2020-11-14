package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func MiddlewareFunc(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		fmt.Println("Before handle request")
		startTime := time.Now()
		handler.ServeHTTP(wr, r)
		fmt.Printf("After handling the request.Response time %s\n", time.Since(startTime))
	})
}
