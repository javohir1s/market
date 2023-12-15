package models

type IncomeProductPrimaryKey struct {
	Id string `json:"id"`
}

type CreateIncomeProduct struct {
	IncomeID    string  `json:"income_id"`
	CategoryID  string  `json:"category_id"`
	ProductName string  `json:"product_name"`
	Barcode     string  `json:"barcode"`
	Quantity    int64   `json:"quantity"`
	IncomePrice float64 `json:"income_price"`
}

type IncomeProduct struct {
	Id          string  `json:"id"`
	IncomeID    string  `json:"income_id"`
	CategoryID  string  `json:"category_id"`
	ProductName string  `json:"product_name"`
	Barcode     string  `json:"barcode"`
	Quantity    int64   `json:"quantity"`
	IncomePrice float64 `json:"income_price"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type UpdateIncomeProduct struct {
	Id          string  `json:"id"`
	ProductName string  `json:"product_name"`
	Barcode     string  `json:"barcode"`
	Quantity    int64   `json:"quantity"`
	IncomePrice float64 `json:"income_price"`
	IncomeID    string  `json:"income_id"`
	CategoryID  string  `json:"category_id"`
}

type GetListIncomeProductRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListIncomeProductResponse struct {
	Count          int              `json:"count"`
	IncomeProducts []*IncomeProduct `json:"income_products"`
}
