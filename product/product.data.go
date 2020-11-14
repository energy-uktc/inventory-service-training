package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var productList = make([]Product, 0)

func init() {
	fileName := "products.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		log.Fatalf("file [%s] does not exist", fileName)
	}

	file, _ := ioutil.ReadFile(fileName)
	err = json.Unmarshal(file, &productList)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d products loaded...\n", len(productList))
}

func getNextProductID() int {
	productID := 0
	for _, product := range productList {
		if productID < product.ProductID {
			productID = product.ProductID
		}
	}
	return productID + 1
}

func getProduct(productID int) (*Product, int) {
	for index, product := range productList {
		if productID == product.ProductID {
			return &product, index
		}
	}
	return nil, 0
}
