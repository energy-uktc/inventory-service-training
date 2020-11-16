package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
)

var productMap = struct {
	sync.RWMutex
	m map[int]Product
}{m: make(map[int]Product)}

func init() {
	var err error
	productMap.m, err = loadProductMap()
	if err != nil {
		log.Fatal(err)
	}
}
func writeProduct(product Product) (*Product, error) {
	addOrUpdateProductID := -1
	if product.ProductID > 0 {
		oldProduct := getProduct(product.ProductID)
		if oldProduct == nil {
			return nil, fmt.Errorf("product id [%d] doesn't exist", product.ProductID)
		}
		addOrUpdateProductID = product.ProductID
	} else {
		addOrUpdateProductID = getNextProductID()
		product.ProductID = addOrUpdateProductID
	}
	productMap.Lock()
	defer productMap.Unlock()
	productMap.m[product.ProductID] = product
	return &product, nil
}

func getNextProductID() int {
	productID := 0
	productMap.RLock()
	for _, product := range productMap.m {
		if productID < product.ProductID {
			productID = product.ProductID
		}
	}
	productMap.RUnlock()
	return productID + 1
}

func getProduct(productID int) *Product {
	productMap.RLock()
	defer productMap.RUnlock()
	if product, ok := productMap.m[productID]; ok {
		return &product
	}
	return nil
}

func getProductList() []Product {
	productMap.RLock()
	defer productMap.RUnlock()
	products := make([]Product, 0, len(productMap.m))
	for _, value := range productMap.m {
		products = append(products, value)
	}
	sort.Slice(products, func(i, j int) bool {
		return products[i].ProductID < products[j].ProductID
	})
	return products
}

func removeProduct(productID int) {
	productMap.Lock()
	defer productMap.Unlock()
	delete(productMap.m, productID)
}

func loadProductMap() (map[int]Product, error) {
	productList := make([]Product, 0)
	fileName := "products.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		log.Fatalf("file [%s] does not exist", fileName)
		return nil, err
	}

	file, _ := ioutil.ReadFile(fileName)
	err = json.Unmarshal(file, &productList)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	prodMap := make(map[int]Product)
	for i := 0; i < len(productList); i++ {
		prodMap[productList[i].ProductID] = productList[i]
	}
	fmt.Printf("%d products loaded...\n", len(productList))
	return prodMap, nil
}
