package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func MiddlewareFunc(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		fmt.Println("Before handle request")
		wr.Header().Add("Content-Type", "application/json")
		wr.Header().Add("Access-Control-Allow-Origin", "*")
		wr.Header().Add("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT,OPTION")
		wr.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		startTime := time.Now()
		handler.ServeHTTP(wr, r)
		fmt.Printf("After handling the request.Response time %s\n", time.Since(startTime))
	})
}
