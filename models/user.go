package models

type UserPrimaryKey struct {
	Id string `json:"id"`
}

type CreateUser struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	Active     bool   `json:"active"`
	ClientType string `json:"client_type"`
}

type User struct {
	Id         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	Active     bool   `json:"active"`
	ClientType string `json:"client_type"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type UpdateUser struct {
	Id         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	Active     bool   `json:"active"`
	ClientType string `json:"client_type"`
}

type GetListUserRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
}

type GetListUserResponse struct {
	Count int64   `json:"count"`
	User  []*User `json:"users"`
}
