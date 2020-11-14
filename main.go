package main

import (
	"log"
	"net/http"
	"time"

	"github.com/energy-uktc/inventory-service-training/product"
)

type Test struct {
	Message string
}

func (t *Test) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	log.Fatal("err")
	time.Sleep(30 * time.Second)
	wr.Write([]byte(t.Message))
}

func simpleFunction(wr http.ResponseWriter, r *http.Request) {
	wr.Write([]byte("Simple function has been called"))
}
func main() {
	http.HandleFunc("/api/simple", simpleFunction)
	http.Handle("/api/test", &Test{Message: "Hello 2"})
	product.Setup("/api")
	http.ListenAndServe(":8081", nil)
}
