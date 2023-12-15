package models

type RemainderPrimaryKey struct {
	Id string `json:"id"`
}

type CreateRemainder struct {
	BranchID    string `json:"branch_id"`
	CategoryID  string `json:"category_id"`
	ProductName string `json:"product_name"`
	Barcode     string `json:"barcode"`
	PriceIncome float64 `json:"price_income"`
	Quantity    int    `json:"quantity"`
}

type Remainder struct {
	Id          string  `json:"id"`
	BranchID    string  `json:"branch_id"`
	CategoryID  string  `json:"category_id"`
	ProductName string  `json:"product_name"`
	Barcode     string  `json:"barcode"`
	PriceIncome float64 `json:"price_income"`
	Quantity    int     `json:"quantity"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type UpdateRemainder struct {
	Id          string  `json:"id"`
	ProductName string  `json:"product_name"`
	Barcode     string  `json:"barcode"`
	PriceIncome float64 `json:"price_income"`
	Quantity    int     `json:"quantity"`
}

type GetListRemainderRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListRemainderResponse struct {
	Count     int         `json:"count"`
	Remainder []*Remainder `json:"remainder"`
}
