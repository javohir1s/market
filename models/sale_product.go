package models

type SaleProductPrimaryKey struct {
	Id string `json:"id"`
}

type CreateSaleProduct struct {
	SaleID            string  `json:"sale_id"`
	CategoryID        string  `json:"category_id"`
	ProductName       string  `json:"product_name"`
	Barcode           string  `json:"barcode"`
	RemainingQuantity int     `json:"remaining_quantity"`
	Quantity          int     `json:"quantity"`
	AllowDiscount     bool    `json:"allow_discount"`
	DiscountType      string  `json:"discount_type"`
	Discount          float64 `json:"discount"`
	Price             float64 `json:"price"`
	TotalAmount       float64 `json:"total_amount"`
}

type SaleProduct struct {
	Id                string  `json:"id"`
	SaleID            string  `json:"sale_id"`
	CategoryID        string  `json:"category_id"`
	ProductName       string  `json:"product_name"`
	Barcode           string  `json:"barcode"`
	RemainingQuantity int     `json:"remaining_quantity"`
	Quantity          int     `json:"quantity"`
	AllowDiscount     bool    `json:"allow_discount"`
	DiscountType      string  `json:"discount_type"`
	Discount          float64 `json:"discount"`
	Price             float64 `json:"price"`
	TotalAmount       float64 `json:"total_amount"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

type UpdateSaleProduct struct {
	Id                string  `json:"id"`
	RemainingQuantity int     `json:"remaining_quantity"`
	Quantity          int     `json:"quantity"`
	AllowDiscount     bool    `json:"allow_discount"`
	DiscountType      string  `json:"discount_type"`
	Discount          float64 `json:"discount"`
	Price             float64 `json:"price"`
	TotalAmount       float64 `json:"total_amount"`
}

type GetListSaleProductRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListSaleProductResponse struct {
	Count        int             `json:"count"`
	SaleProducts []*SaleProduct `json:"sale_products"`
}
