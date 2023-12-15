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

type incomeRepo struct {
	db *pgxpool.Pool
}

func NewIncomeRepo(db *pgxpool.Pool) *incomeRepo {
	return &incomeRepo{
		db: db,
	}
}

func (r *incomeRepo) Create(ctx context.Context, req *models.CreateIncome) (*models.Income, error) {

	var (
		incomeId = uuid.New().String()
		query    = `
			INSERT INTO income(
				id,
				branch_id,
				supplier_id,
				date_time,
				status,
				updated_at
			) VALUES ($1, $2, $3, $4, $5, NOW())`
	)

	_, err := r.db.Exec(ctx,
		query,
		incomeId,
		helpers.NewNullString(req.BranchID),
		helpers.NewNullString(req.SupplierID),
		req.DateTime,
		req.Status,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.IncomePrimaryKey{Id: incomeId})
}

func (r *incomeRepo) GetByID(ctx context.Context, req *models.IncomePrimaryKey) (*models.Income, error) {

	var (
		query = `
			SELECT
				 id,
				 branch_id,
				 supplier_id,
				 date_time,
				 status,
				 created_at,
				 updated_at
			FROM income
			WHERE id = $1
		`
	)

	var (
		Id         sql.NullString
		BranchID   sql.NullString
		SupplierID sql.NullString
		DateTime   sql.NullString
		Status     sql.NullString
		CreatedAt  sql.NullString
		UpdatedAt  sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&Id,
		&BranchID,
		&SupplierID,
		&DateTime,
		&Status,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Income{
		Id:         Id.String,
		BranchID:   BranchID.String,
		SupplierID: SupplierID.String,
		DateTime:   DateTime.String,
		Status:     Status.String,
		CreatedAt:  CreatedAt.String,
		UpdatedAt:  UpdatedAt.String,
	}, nil
}

func (r *incomeRepo) GetList(ctx context.Context, req *models.GetListIncomeRequest) (*models.GetListIncomeResponse, error) {
	var (
		resp   models.GetListIncomeResponse
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
		where += " AND (branch_id ILIKE '%" + req.Search + "%' OR category_id ILIKE '%" + req.Search + "%' OR barcode ILIKE '%" + req.Search + "%' OR income_id ILIKE '%" + req.Search + "%')"
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			id,
			branch_id,
			supplier_id,
			date_time,
			status,
			created_at,
			updated_at
		FROM income
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			Id         sql.NullString
			BranchID   sql.NullString
			SupplierID sql.NullString
			DateTime   sql.NullString
			Status     sql.NullString
			CreatedAt  sql.NullString
			UpdatedAt  sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&Id,
			&BranchID,
			&SupplierID,
			&DateTime,
			&Status,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		resp.Incomes = append(resp.Incomes, &models.Income{
			Id:         Id.String,
			BranchID:   BranchID.String,
			SupplierID: SupplierID.String,
			DateTime:   DateTime.String,
			Status:     Status.String,
			CreatedAt:  CreatedAt.String,
			UpdatedAt:  UpdatedAt.String,
		})
	}

	return &resp, nil
}

func (r *incomeRepo) Update(ctx context.Context, req *models.UpdateIncome) (int64, error) {

	query := `
		UPDATE income
			SET
				branch_id = $2,
				supplier_id = $3,
				date_time = $4,
				status = $5, 
				updated_at = NOW()
		WHERE id = $1
	`
	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		helpers.NewNullString(req.BranchID),
		helpers.NewNullString(req.SupplierID),
		req.DateTime,
		req.Status,
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *incomeRepo) Delete(ctx context.Context, req *models.IncomePrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM income WHERE id = $1", req.Id)
	return err
}
