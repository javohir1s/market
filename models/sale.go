package models

type SalePrimaryKey struct {
	Id string `json:"id"`
}

type CreateSale struct {
	SaleID      string `json:"sale_id"`
	BranchID    string `json:"branch_id"`
	SalePointID string `json:"salepoint_id"`
	ShiftID     string `json:"shift_id"`
	EmployeeID  string `json:"employee_id"`
	Barcode     string `json:"barcode"`
	Status      string `json:"status"`
}

type Sale struct {
	Id          string `json:"id"`
	SaleID      string `json:"sale_id"`
	BranchID    string `json:"branch_id"`
	SalePointID string `json:"salepoint_id"`
	ShiftID     string `json:"shift_id"`
	EmployeeID  string `json:"employee_id"`
	Barcode     string `json:"barcode"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UpdateSale struct {
	Id          string `json:"id"`
	BranchID    string `json:"branch_id"`
	SalePointID string `json:"salepoint_id"`
	ShiftID     string `json:"shift_id"`
	EmployeeID  string `json:"employee_id"`
	Barcode     string `json:"barcode"`
	Status      string `json:"status"`
}

type GetListSaleRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListSaleResponse struct {
	Count int     `json:"count"`
	Sales []*Sale `json:"sales"`
}
