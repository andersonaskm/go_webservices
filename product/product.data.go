package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/andersonaskm/go_webservices/database"
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
func getProduct(productID int) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	row := database.DbConn.QueryRowContext(ctx, `SELECT productId, manufacturer, sku, upc, unitPrice, quantity, productName 
		FROM products WHERE productId = ?`, productID)

	product := &Product{}
	err := row.Scan(&product.ProductID,
		&product.Manufacturer,
		&product.Sku,
		&product.Upc,
		&product.UnitPrice,
		&product.Quantity,
		&product.ProductName)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return product, nil

	// productMap.RLock() // read lock to prevent another thread
	// defer productMap.RUnlock()
	// if product, ok := productMap.m[productID]; ok {
	// 	return &product
	// }
	// return nil
}

// exclui um produto por ID
func removeProduct(productID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM products WHERE productId = ?`, productID)
	if err != nil {
		return err
	}
	return nil
	// productMap.Lock()
	// defer productMap.Unlock()
	// delete(productMap.m, productID)
}

// obtem a lista de produtos
func getProductList() ([]Product, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	results, errQuery := database.DbConn.QueryContext(ctx, `
		SELECT productId, manufacturer, sku, upc, unitPrice, quantity, productName 
		FROM products`)

	if errQuery != nil {
		return nil, errQuery
	}

	defer results.Close()

	products := make([]Product, 0)

	for results.Next() {
		var product Product
		results.Scan(&product.ProductID,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.UnitPrice,
			&product.Quantity,
			&product.ProductName)

		products = append(products, product)
	}

	return products, nil

	// productMap.RLock()
	// products := make([]Product, 0, len(productMap.m))
	// for _, value := range productMap.m {
	// 	products = append(products, value)
	// }
	// productMap.RUnlock()
	// return products
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

// atualiza um produto
func updateProduct(product Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `UPDATE products SET manufacturer=? WHERE productId=?`, product.Manufacturer, product.ProductID)
	if err != nil {
		return err
	}
	return nil
}

func insertProduct(product Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := database.DbConn.ExecContext(ctx, `INSERT INTO products(manufacturer, sku, upc, unitPrice, quantity, productName) 
	VALUES (?, ?, ?, ?, ?, ?)`, product.Manufacturer, product.Sku, product.Upc, product.UnitPrice, product.Quantity, product.ProductName)
	if err != nil {
		return 0, nil
	}

	insertID, errLastInsertId := result.LastInsertId()
	if errLastInsertId != nil {
		return 0, nil
	}

	return int(insertID), nil
}

// incluir ou atualizar um produto
func addOrUpdateProduct(product Product) (int, error) {
	addOrUpdateID := -1

	if product.ProductID > 0 {

		oldProduct, errGetProduct := getProduct(product.ProductID)

		if errGetProduct != nil {
			return addOrUpdateID, errGetProduct
		}

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
