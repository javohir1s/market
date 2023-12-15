package models

type ProductPrimaryKey struct {
	Id string `json:"id"`
}

type CreateProduct struct {
	Photo      string  `json:"photo"`
	Title      string  `json:"title"`
	CategoryID string  `json:"category_id"`
	Barcode    string  `json:"barcode"`
	Price      float64 `json:"price"`
}

type Product struct {
	Id         string  `json:"id"`
	Photo      string  `json:"photo"`
	Title      string  `json:"title"`
	CategoryID string  `json:"category_id"`
	Barcode    string  `json:"barcode"`
	Price      float64 `json:"price"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type UpdateProduct struct {
	Id         string  `json:"id"`
	Photo      string  `json:"photo"`
	Title      string  `json:"title"`
	CategoryID string  `json:"category_id"`
	Barcode    string  `json:"barcode"`
	Price      float64 `json:"price"`
}

type GetListProductRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListProductResponse struct {
	Count    int        `json:"count"`
	Products []*Product `json:"products"`
}
