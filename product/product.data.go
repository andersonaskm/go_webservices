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

// inicializa o package
func init() {
	fmt.Println("Loading products...")
	prodMap, err := loadProductMap()
	productMap.m = prodMap
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d products loaded.... \n", len(productMap.m))
}

// cria um mapa de produtos a partir do arquivo
func loadProductMap() (map[int]Product, error) {
	fileName := "products.json"

	// verify if file exists
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}

	// load file content
	file, _ := ioutil.ReadFile(fileName)
	productList := make([]Product, 0)
	errProductList := json.Unmarshal([]byte(file), &productList)
	if errProductList != nil {
		log.Fatal(errProductList)
	}

	// fill products map
	prodMap := make(map[int]Product)
	for i := 0; i < len(productList); i++ {
		prodMap[productList[i].ProductID] = productList[i]
	}

	return prodMap, nil
}

// obtem um produto pelo ID
func getProduct(productID int) *Product {
	productMap.RLock() // read lock to prevent another thread
	defer productMap.RUnlock()
	if product, ok := productMap.m[productID]; ok {
		return &product
	}
	return nil
}

// exclui um produto por ID
func removeProduct(productID int) {
	productMap.Lock()
	defer productMap.Unlock()
	delete(productMap.m, productID)
}

// obtem a lista de produtos
func getProductList() []Product {
	productMap.RLock()
	products := make([]Product, 0, len(productMap.m))
	for _, value := range productMap.m {
		products = append(products, value)
	}
	productMap.RUnlock()
	return products
}

// obtem os identificadores de produtos ordenados
func getProductIds() []int {
	productMap.RLock()
	productsIds := []int{}
	for key := range productMap.m {
		productsIds = append(productsIds, key)
	}
	productMap.RUnlock()
	sort.Ints(productsIds)
	return productsIds
}

// obtem o novo identificador de produto
func getNextProductID() int {
	productIDs := getProductIds()
	return productIDs[len(productIDs)-1] + 1
}

// incluir ou atualizar um produto
func addOrUpdateProduct(product Product) (int, error) {
	addOrUpdateID := -1

	if product.ProductID > 0 {

		oldProduct := getProduct(product.ProductID)
		if oldProduct == nil {
			return 0, fmt.Errorf("product id [%d] doesn't exist", product.ProductID)
		}
		addOrUpdateID = product.ProductID
	} else {
		addOrUpdateID = getNextProductID()
		product.ProductID = addOrUpdateID
	}

	productMap.Lock()
	productMap.m[addOrUpdateID] = product
	productMap.Unlock()

	return addOrUpdateID, nil
}
