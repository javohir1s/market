package models

type BranchPrimaryKey struct {
	Id string `json:"id"`
}

type CreateBranch struct {
	BranchCode string `json:"branch_code"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
}

type Branch struct {
	Id         string `json:"id"`
	BranchCode string `json:"branch_code"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UpdateBranch struct {
	Id         string `json:"id"`
	BranchCode string `json:"branch_code"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
}

type GetListBranchRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListBranchResponse struct {
	Count    int       `json:"count"`
	Branches []*Branch `json:"branches"`
}
