package models

type IncomePrimaryKey struct {
	Id string `json:"id"`
}

type CreateIncome struct {
	BranchID   string `json:"branch_id"`
	SupplierID string `json:"supplier_id"`
	DateTime   string `json:"date_time"`
	Status     string `json:"status"`
}

type Income struct {
	Id         string `json:"id"`
	BranchID   string `json:"branch_id"`
	SupplierID string `json:"supplier_id"`
	DateTime   string `json:"date_time"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type UpdateIncome struct {
	Id         string `json:"id"`
	BranchID   string `json:"branch_id"`
	SupplierID string `json:"supplier_id"`
	DateTime   string `json:"date_time"`
	Status     string `json:"status"`
}

type GetListIncomeRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListIncomeResponse struct {
	Count   int       `json:"count"`
	Incomes []*Income `json:"incomes"`
}
