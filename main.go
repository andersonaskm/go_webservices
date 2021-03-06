package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/andersonaskm/go_webservices/database"
	"github.com/andersonaskm/go_webservices/product"
	"github.com/andersonaskm/go_webservices/receipt"
	_ "github.com/go-sql-driver/mysql"
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

type BlogPost struct {
	Title   string
	Content string
}

func init() {

}

func (handler *ordersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(handler.Message))
}

func orderDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Order details called!"))
}

const apiBasePath = "/api"

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
	fmt.Println("Connect To DataBase")
	fmt.Println("-----------------------------")
	database.SetUpDatabase()

	fmt.Println("-----------------------------")
	fmt.Println("Template")
	fmt.Println("-----------------------------")
	post := BlogPost{Title: "First Post", Content: "This is the first post content"}
	tmpl, errTemplate := template.New("blog-tmpl").Parse(`<h1>{{.Title}}</h1><div><p>{{.Content}}</p></div>`)
	if errTemplate != nil {
		panic(errTemplate)
	}
	errTemplateExec := tmpl.Execute(os.Stdout, post)
	if errTemplateExec != nil {
		panic(errTemplateExec)
	}

	fmt.Println("-----------------------------")
	fmt.Println("Basic Handler")
	fmt.Println("-----------------------------")

	/*
		func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error
	*/

	product.SetupRoutes(apiBasePath) // products
	receipt.SetupRoutes(apiBasePath) // receipts

	http.Handle("/orders", &ordersHandler{Message: "Orders called!"})
	http.HandleFunc("/orders/details", orderDetailsHandler)

	errHttp := http.ListenAndServe(":5002", nil)
	if errHttp != nil {
		log.Fatal(errHttp)
	}

}
