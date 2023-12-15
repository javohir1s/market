package models

type EmployeePrimaryKey struct {
	Id string `json:"id"`
}

type CreateEmployee struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	BranchID    string `json:"branch_id"`
	SalepointID string `json:"salepoint_id"`
	UserType    string `json:"user_type"`
}

type Employee struct {
	Id         string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	BranchID    string `json:"branch_id"`
	SalepointID string `json:"salepoint_id"`
	UserType    string `json:"user_type"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UpdateEmployee struct {
	Id          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	BranchID    string `json:"branch_id"`
	SalepointID string `json:"salepoint_id"`
	UserType    string `json:"user_type"`
}

type GetListEmployeeRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListEmployeeResponse struct {
	Count     int         `json:"count"`
	Employees []*Employee `json:"employees"`
}
