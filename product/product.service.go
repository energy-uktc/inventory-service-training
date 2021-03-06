package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/energy-uktc/inventory-service-training/middleware"
)

const productsPath = "products"

func handleProduct(wr http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", productsPath))
	if len(urlPathSegments[1:]) > 1 {
		wr.WriteHeader(http.StatusBadRequest)
		return
	}
	switch r.Method {

	case http.MethodGet:
		productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
		if err != nil {
			log.Print(err)
			wr.WriteHeader((http.StatusNotFound))
			return
		}
		product := getProduct(productID)
		if product == nil {
			wr.WriteHeader(http.StatusNotFound)
			return
		}
		productJSON, err := json.MarshalIndent(product, "", "	")
		if err != nil {
			log.Print(err)
			wr.WriteHeader(http.StatusInternalServerError)
			return
		}
		wr.Header().Set("Content-Type", "application/json")
		wr.Write(productJSON)

	case http.MethodPut:
		productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
		if err != nil {
			log.Print(err)
			wr.WriteHeader((http.StatusNotFound))
			return
		}
		product := getProduct(productID)
		if product == nil {
			wr.WriteHeader(http.StatusNotFound)
			return
		}
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			wr.WriteHeader(http.StatusBadRequest)
			return
		}
		var updatedProduct Product
		err = json.Unmarshal(reqBody, &updatedProduct)
		if err != nil {
			log.Println(err)
			wr.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = writeProduct(updatedProduct)
		if err != nil {
			wr.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
		if err != nil {
			log.Print(err)
			wr.WriteHeader((http.StatusNotFound))
			return
		}
		product := getProduct(productID)
		if product == nil {
			wr.WriteHeader(http.StatusNotFound)
			return
		}
		removeProduct(productID)
	default:
		wr.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func handleProducts(wr http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		products, err := json.MarshalIndent(getProductList(), "", "	")
		if err != nil {
			log.Fatal(err)
		}
		wr.Header().Set("Content-Type", "application/json")
		_, err = wr.Write(products)
		if err != nil {
			log.Fatal(err)
		}

	case http.MethodPost:
		var product Product
		err := json.NewDecoder(r.Body).Decode(&product)

		if err != nil {
			log.Print(err)
			wr.WriteHeader(http.StatusBadRequest)
			return
		}
		if product.ProductID != 0 {
			log.Printf("ProductID should be [0] current value %v", product.ProductID)
			wr.WriteHeader(http.StatusBadRequest)
			return
		}
		createdProduct, err := writeProduct(product)
		if err != nil {
			wr.WriteHeader(http.StatusBadRequest)
			return
		}
		p, _ := json.MarshalIndent(createdProduct, "", "	")
		wr.Header().Set("Content-Type", "application/json")
		wr.WriteHeader(http.StatusCreated)
		wr.Write(p)

	case http.MethodOptions:
		return
	default:
		wr.WriteHeader(http.StatusMethodNotAllowed)

	}

}

//Setup products route
func SetupRoutes(apiBasePath string) {
	productsHandler := http.HandlerFunc(handleProducts)
	productHandler := http.HandlerFunc(handleProduct)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productsPath), middleware.MiddlewareFunc(productsHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productsPath), middleware.MiddlewareFunc(productHandler))
}
