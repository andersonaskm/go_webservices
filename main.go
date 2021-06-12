package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	http.ListenAndServe(":5002", nil)

}
