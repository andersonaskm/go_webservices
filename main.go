package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ordersHandler struct {
	Message string
}

type Order struct {
	OrderID    int     `json:"orderId"`
	CustomerId int     `json:"customerId"`
	TotalValue float64 `json:"totalValue"`
}

type OrderDetail struct {
	ProductID int     `json:"productId"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unitPrice"`
}

var productList []Product

func init() {
	productsJSON := `[
		{
			"productId": 1,
			"manufactures": "Johns-Jenkins",
			"sku": "p5z343vdS",
			"upc": "939581000000",
			"unitPrice": "497.45",
			"quantity": 9703,
			"productName": "sticky note"
		},
		{
			"productId": 2,
			"manufactures": "Hessel, SChimmel and Feeney",
			"sku": "i7v300kmx",
			"upc": "740979000000",
			"unitPrice": "282.29",
			"quantity": 9217,
			"productName": "leg warmers"
		},
		{
			"productId": 3,
			"manufactures": "Swaniawski, Bartoletti and Bruen",
			"sku": "q0L657ys7",
			"upc": "111730000000",
			"unitPrice": "436.26",
			"quantity": 5905,
			"productName": "lamp shade"
		}
	]`

	err := json.Unmarshal([]byte(productsJSON), &productList)
	if err != nil {
		log.Fatal(err)
	}
}

func getNextID() int {
	highestID := -1
	for _, product := range productList {
		if highestID < product.ProductID {
			highestID = product.ProductID
		}
	}

	return highestID + 1
}

func findByProductID(productID int) (*Product, int) {
	for i, product := range productList {
		if product.ProductID == productID {
			return &product, i
		}
	}

	return nil, 0
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productsJson, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(productsJson)
	case http.MethodPost:
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		errProduct := json.Unmarshal(bodyBytes, &newProduct)
		if errProduct != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		if newProduct.ProductID != 0 {
			w.WriteHeader(http.StatusBadRequest)
		}

		newProduct.ProductID = getNextID()
		productList = append(productList, newProduct)
		w.WriteHeader(http.StatusCreated)
		return
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "products/")
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	product, listIndex := findByProductID(productID)
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		productJSON, err := json.Marshal(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productJSON)
	case http.MethodPut:
		var updatedProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		errUpdatedProduct := json.Unmarshal(bodyBytes, &updatedProduct)
		if errUpdatedProduct != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updatedProduct.ProductID != productID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		product = &updatedProduct
		productList[listIndex] = *product
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (handler *ordersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(handler.Message))
}

func orderDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Order details called!"))
}

func main() {

	fmt.Println("-----------------------------")
	fmt.Println("JSON Marshal")
	fmt.Println("-----------------------------")
	order := &Order{
		OrderID:    12,
		CustomerId: 1,
		TotalValue: 45.50,
	}

	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(orderJSON))

	fmt.Println("-----------------------------")
	fmt.Println("JSON Unmarshal")
	fmt.Println("-----------------------------")

	orderDetailJSON := `{
		"productId": 10,
		"quantity": 1,
		"unitPrice": 15.50
	}`

	orderDetail := OrderDetail{}
	errOrderDetail := json.Unmarshal([]byte(orderDetailJSON), &orderDetail)
	if errOrderDetail != nil {
		log.Fatal(errOrderDetail)
	}

	fmt.Println(orderDetail.ProductID, orderDetail.Quantity, orderDetail.UnitPrice)

	fmt.Println("-----------------------------")
	fmt.Println("Basic Handler")
	fmt.Println("-----------------------------")

	http.Handle("/orders", &ordersHandler{Message: "Orders called!"})
	http.HandleFunc("/orders/details", orderDetailsHandler)
	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/products/", productHandler)
	http.ListenAndServe(":5002", nil)

}
