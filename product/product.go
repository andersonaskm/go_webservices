package product

type Product struct {
	ProductID    int     `json:"productId"`
	Manufacturer string  `json:"manufacturer"`
	Sku          string  `json:"sku"`
	Upc          string  `json:"upc"`
	UnitPrice    float32 `json:"unitPrice"`
	Quantity     int     `json:"quantity"`
	ProductName  string  `json:"productName"`
}
