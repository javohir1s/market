package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"market_system/models"
	"market_system/pkg/helpers"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type employeeRepo struct {
	db *pgxpool.Pool
}

func NewEmployeeRepo(db *pgxpool.Pool) *employeeRepo {
	return &employeeRepo{
		db: db,
	}
}

func (r *employeeRepo) Create(ctx context.Context, req *models.CreateEmployee) (*models.Employee, error) {

	var (
		employeeID = uuid.New().String()
		query      = `
			INSERT INTO employee(
				id,
				first_name, 
				last_name, 
				phone, 
				login, 
				password, 
				branch_id, 
				salepoint_id, 
				user_type, 
				updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		`
	)

	_, err := r.db.Exec(ctx,
		query,
		employeeID,
		req.FirstName,
		req.LastName,
		req.Phone,
		req.Login,
		req.Password,
		helpers.NewNullString(req.BranchID),
		helpers.NewNullString(req.SalepointID),
		req.UserType,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.EmployeePrimaryKey{Id: employeeID})
}

func (r *employeeRepo) GetByID(ctx context.Context, req *models.EmployeePrimaryKey) (*models.Employee, error) {

	var (
		query = `
			SELECT
				id,
				first_name,
				last_name,
				phone,
				login,
				password,
				branch_id,
				salepoint_id,
				user_type,
				created_at,
				updated_at
			FROM employee
			WHERE id = $1
		`
	)

	var (
		ID          sql.NullString
		FirstName   sql.NullString
		LastName    sql.NullString
		Phone       sql.NullString
		Login       sql.NullString
		Password    sql.NullString
		BranchID    sql.NullString
		SalepointID sql.NullString
		UserType    sql.NullString
		CreatedAt   sql.NullString
		UpdatedAt   sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&ID,
		&FirstName,
		&LastName,
		&Phone,
		&Login,
		&Password,
		&BranchID,
		&SalepointID,
		&UserType,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Employee{
		Id:          ID.String,
		FirstName:   FirstName.String,
		LastName:    LastName.String,
		Phone:       Phone.String,
		Login:       Login.String,
		Password:    Password.String,
		BranchID:    BranchID.String,
		SalepointID: SalepointID.String,
		UserType:    UserType.String,
		CreatedAt:   CreatedAt.String,
		UpdatedAt:   UpdatedAt.String,
	}, nil
}


func (r *employeeRepo) GetList(ctx context.Context, req *models.GetListEmployeeRequest) (*models.GetListEmployeeResponse, error) {
	var (
		resp   models.GetListEmployeeResponse
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
		sort   = " ORDER BY created_at DESC"
	)

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if len(req.Search) > 0 {
		where += " AND (first_name ILIKE '%" + req.Search + "%' OR last_name ILIKE '%" + req.Search + "%' OR phone ILIKE '%" + req.Search + "%')"
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			id,
			first_name,
			last_name,
			phone,
			login,
			branch_id,
			salepoint_id,
			user_type,
			created_at,
			updated_at
		FROM employee
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			ID          sql.NullString
			FirstName   sql.NullString
			LastName    sql.NullString
			Phone       sql.NullString
			Login       sql.NullString
			Password    sql.NullString
			BranchID    sql.NullString
			SalepointID sql.NullString
			UserType    sql.NullString
			CreatedAt   sql.NullString
			UpdatedAt   sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&ID,
			&FirstName,
			&LastName,
			&Phone,
			&Login,
			&Password,
			&BranchID,
			&SalepointID,
			&UserType,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Employees = append(resp.Employees, &models.Employee{
			Id:          ID.String,
			FirstName:   FirstName.String,
			LastName:    LastName.String,
			Phone:       Phone.String,
			Login:       Login.String,
			Password:    Password.String,
			BranchID:    BranchID.String,
			SalepointID: SalepointID.String,
			UserType:    UserType.String,
			CreatedAt:   CreatedAt.String,
			UpdatedAt:   UpdatedAt.String,
		})
	}

	return &resp, nil
}

func (r *employeeRepo) Update(ctx context.Context, req *models.UpdateEmployee) (int64, error) {

	query := `
		UPDATE employee
			SET
				first_name = $2,
				last_name = $3,
				phone = $4,
				login = $5,
				password = $6,
				branch_id = $7,
				salepoint_id = $8,
				user_type = $9, 
				updated_at = NOW()
		WHERE id = $1
	`
	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		req.FirstName,
		req.LastName,
		req.Phone,
		req.Login,
		req.Password,
		req.BranchID,
		req.SalepointID,
		req.UserType,
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *employeeRepo) Delete(ctx context.Context, req *models.EmployeePrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM employee WHERE id = $1", req.Id)
	return err
}
