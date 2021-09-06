package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/andersonaskm/go_webservices/cors"
	"golang.org/x/net/websocket"
)

const productsBasePath = "products"

// configura as rotas de navegação
func SetupRoutes(apiBasePath string) {
	handleProducts := http.HandlerFunc(productsHandler)
	handleProduct := http.HandlerFunc(productHandler)
	handleReports := http.HandlerFunc(productReportHandler)

	http.Handle("/websocket", websocket.Handler(productWebSocket))
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productsBasePath), cors.Middleware(handleProducts))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productsBasePath), cors.Middleware(handleProduct))
	http.Handle(fmt.Sprintf("%s/%s/reports", apiBasePath, productsBasePath), cors.Middleware(handleReports))
}

// handler da rota '/products'
func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		// call getProductList from Data
		productList, errGetProductList := getProductList()
		if errGetProductList != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

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

		// call addOrUddateProduct from Data
		_, errInsertProduct := insertProduct(newProduct)
		if errInsertProduct != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return

	case http.MethodOptions:
		return
	}
}

// handler da rota '/products/'
func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "products/")
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// call getProduct from Data
	product, errGetProduct := getProduct(productID)
	if errGetProduct != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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

		// call addOrUpdateProduct from Data
		errUpdateProduct := updateProduct(updatedProduct)
		if errUpdateProduct != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//addOrUpdateProduct(updatedProduct)
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		removeProduct(productID)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
