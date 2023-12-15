package models

type TransactionPrimaryKey struct {
	Id string `json:"id"`
}

type CreateTransaction struct {
	ShiftID     string  `json:"sale_id"`
	Cash        float64 `json:"cash"`
	Uzcard      float64 `json:"uzcard"`
	Payme       float64 `json:"payme"`
	Click       float64 `json:"click"`
	Humo        float64 `json:"humo"`
	Apelsin     float64 `json:"apelsin"`
	TotalAmount float64 `json:"total_amount"`
}

type Transaction struct {
	Id          string  `json:"id"`
	ShiftID     string  `json:"shift_id"`
	Cash        float64 `json:"cash"`
	Uzcard      float64 `json:"uzcard"`
	Payme       float64 `json:"payme"`
	Click       float64 `json:"click"`
	Humo        float64 `json:"humo"`
	Apelsin     float64 `json:"apelsin"`
	TotalAmount float64 `json:"total_amount"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type UpdateTransaction struct {
	Id          string  `json:"id"`
	Cash        float64 `json:"cash"`
	Uzcard      float64 `json:"uzcard"`
	Payme       float64 `json:"payme"`
	Click       float64 `json:"click"`
	Humo        float64 `json:"humo"`
	Apelsin     float64 `json:"apelsin"`
	TotalAmount float64 `json:"total_amount"`
}

type GetListTransactonRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListTransactionResponse struct {
	Count        int            `json:"count"`
	Transactions []*Transaction `json:"transactions"`
}
