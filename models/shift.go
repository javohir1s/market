package models

type ShiftPrimaryKey struct {
	Id string `json:"id"`
}

type CreateShift struct {
	BranchID    string `json:"branch_id"`
	EmployeeID  string `json:"employee_id"`
	SalePointID string `json:"sale_point_id"`
	Status      string `json:"status"`
	OpenShift   string `json:"open_shift"`
	CloseShift  string `json:"close_shift"`
}

type Shift struct {
	Id          string `json:"id"`
	BranchID    string `json:"branch_id"`
	EmployeeID  string `json:"employee_id"`
	SalePointID string `json:"sale_point_id"`
	Status      string `json:"status"`
	OpenShift   string `json:"open_shift"`
	CloseShift  string `json:"close_shift"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UpdateShift struct {
	Id          string `json:"id"`
	BranchID    string `json:"branch_id"`
	EmployeeID  string `json:"employee_id"`
	SalePointID string `json:"sale_point_id"`
	Status      string `json:"status"`
	OpenShift   string `json:"open_shift"`
	CloseShift  string `json:"close_shift"`
}

type GetListShiftRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListShiftResponse struct {
	Count int      `json:"count"`
	Shift []*Shift `json:"shift"`
}
