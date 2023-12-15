package models

type SupplierPrimaryKey struct {
	Id string `json:"id"`
}

type CreateSupplier struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	IsActive    bool   `json:"is_active"`
}

type Supplier struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UpdateSupplier struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	IsActive    bool   `json:"is_active"`
}

type GetListSupplierRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListSupplierResponse struct {
	Count     int         `json:"count"`
	Suppliers []*Supplier `json:"suppliers"`
}
