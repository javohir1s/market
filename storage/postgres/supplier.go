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

type supplierRepo struct {
	db *pgxpool.Pool
}

func NewSupplierRepo(db *pgxpool.Pool) *supplierRepo {
	return &supplierRepo{
		db: db,
	}
}

func (r *supplierRepo) Create(ctx context.Context, req *models.CreateSupplier) (*models.Supplier, error) {

	var (
		supplierID = uuid.New().String()
		query      = `
			INSERT INTO supplier(
				id,
				name, 
				phone_number, 
				is_active,
				updated_at
			) VALUES ($1, $2, $3, $4, NOW())`
	)

	_, err := r.db.Exec(ctx,
		query,
		supplierID,
		req.Name,
		helpers.NewNullString(req.PhoneNumber),
		req.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.SupplierPrimaryKey{Id: supplierID})
}

func (r *supplierRepo) GetByID(ctx context.Context, req *models.SupplierPrimaryKey) (*models.Supplier, error) {

	var (
		query = `
			SELECT
				id,
				name,
				COALESCE(CAST(phone_number AS VARCHAR), ''),
				is_active,
				created_at,
				updated_at	
			FROM  supplier
			WHERE id = $1
		`
	)

	var (
		ID          sql.NullString
		Name        sql.NullString
		PhoneNumber sql.NullString
		IsActive    sql.NullBool
		CreatedAt   sql.NullString
		UpdatedAt   sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&ID,
		&Name,
		&PhoneNumber,
		&IsActive,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Supplier{
		Id:          ID.String,
		Name:        Name.String,
		PhoneNumber: PhoneNumber.String,
		IsActive:    IsActive.Bool,
		CreatedAt:   CreatedAt.String,
		UpdatedAt:   UpdatedAt.String,
	}, nil
}

func (r *supplierRepo) GetList(ctx context.Context, req *models.GetListSupplierRequest) (*models.GetListSupplierResponse, error) {
	var (
		resp   models.GetListSupplierResponse
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
		where += " AND (name ILIKE '%" + req.Search + "%' OR phone_number ILIKE '%" + req.Search + "%')"
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			id,
			name,
			phone_number,
			is_active,
			created_at,
			updated_at
		FROM supplier
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			Id          sql.NullString
			Name        sql.NullString
			PhoneNumber sql.NullString
			IsActive    sql.NullBool
			CreatedAt   sql.NullString
			UpdatedAt   sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&Id,
			&Name,
			&PhoneNumber,
			&IsActive,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Suppliers = append(resp.Suppliers, &models.Supplier{
			Id:          Id.String,
			Name:        Name.String,
			PhoneNumber: PhoneNumber.String,
			IsActive:    IsActive.Bool,
			CreatedAt:   CreatedAt.String,
			UpdatedAt:   UpdatedAt.String,
		})
	}

	return &resp, nil
}

func (r *supplierRepo) Update(ctx context.Context, req *models.UpdateSupplier) (int64, error) {

	query := `
		UPDATE supplier
			SET
				name = $2,
				phone_number = $3,
				is_active = $4,
				updated_at = NOW()
		WHERE id = $1
	`
	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		req.Name,
		helpers.NewNullString(req.PhoneNumber),
		req.IsActive,
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *supplierRepo) Delete(ctx context.Context, req *models.SupplierPrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM supplier WHERE id = $1", req.Id)
	return err
}
