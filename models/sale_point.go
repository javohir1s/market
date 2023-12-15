package models

type SalePointPrimaryKey struct {
	Id string `json:"id"`
}

type CreateSalePoint struct {
	Branch_id string `json:"branch_id"`
	Name      string `json:"name"`
}

type SalePoint struct {
	Id        string `json:"id"`
	Branch_id string `json:"branch_id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UpdateSalePoint struct {
	Id        string `json:"id"`
	Branch_id string `json:"branch_id"`
	Name      string `json:"name"`
}

type GetListSalePointRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListSalePointResponse struct {
	Count      int          `json:"count"`
	SalePoints []*SalePoint `json:"sale_points"`
}
