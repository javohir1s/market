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

type saleRepo struct {
	db *pgxpool.Pool
}

func NewSaleRepo(db *pgxpool.Pool) *saleRepo {
	return &saleRepo{
		db: db,
	}
}

func (r *saleRepo) Create(ctx context.Context, req *models.CreateSale) (*models.Sale, error) {

	var (
		saleId = uuid.New().String()
		query  = `
			INSERT INTO sale(
				id,
				sale_id,
				branch_id,
				salepoint_id,
				shift_id,
				employee_id,
				barcode,
				status,
				updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW()) 
		`
	)

	_, err := r.db.Exec(ctx,
		query,
		saleId,
		req.SaleID,
		req.BranchID,
		req.SalePointID,
		req.ShiftID,
		req.EmployeeID,
		helpers.NewNullString(req.Barcode),
		helpers.NewNullString(req.Status),
	)

	if err != nil {
		return nil, err

	}
	return r.GetByID(ctx, &models.SalePrimaryKey{Id: saleId})
}

func (r *saleRepo) GetByID(ctx context.Context, req *models.SalePrimaryKey) (*models.Sale, error) {

	var (
		query = `
			SELECT
				id,
				sale_id,
				branch_id,
				salepoint_id,
				shift_id,
				employee_id,
				barcode,
				status,
				created_at,
				updated_at
			FROM  sale
			WHERE id = $1
		`
	)

	var (
		id          sql.NullString
		saleId      sql.NullString
		branchId    sql.NullString
		salepointId sql.NullString
		shiftId     sql.NullString
		employeeId  sql.NullString
		barcode     sql.NullString
		status      sql.NullString
		createdAt   sql.NullString
		updatedAt   sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&saleId,
		&branchId,
		&salepointId,
		&shiftId,
		&employeeId,
		&barcode,
		&status,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Sale{
		Id:          id.String,
		SaleID:      saleId.String,
		BranchID:    branchId.String,
		SalePointID: salepointId.String,
		ShiftID:     shiftId.String,
		EmployeeID:  employeeId.String,
		Barcode:     barcode.String,
		Status:      status.String,
		CreatedAt:   createdAt.String,
		UpdatedAt:   updatedAt.String,
	}, nil
}

func (r *saleRepo) GetList(ctx context.Context, req *models.GetListSaleRequest) (*models.GetListSaleResponse, error) {
	var (
		resp   models.GetListSaleResponse
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
		where += " AND (sale_id ILIKE '%" + req.Search + "%' OR shift_id ILIKE '%" + req.Search + "%' OR branch_id ILIKE '%" + req.Search + "%' OR category_id ILIKE '%" + req.Search + "%' OR barcode ILIKE '%" + req.Search + "%' OR employee_id ILIKE '%" + req.Search + "%')"
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			id,
			sale_id,
			branch_id,
			salepoint_id,
			shift_id,
			employee_id,
			barcode,
			status,
			created_at,
			updated_at
		FROM sale
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id          sql.NullString
			saleId      sql.NullString
			branchId    sql.NullString
			salepointId sql.NullString
			shiftId     sql.NullString
			employeeId  sql.NullString
			barcode     sql.NullString
			status      sql.NullString
			createdAt   sql.NullString
			updatedAt   sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&id,
			&saleId,
			&branchId,
			&salepointId,
			&shiftId,
			&employeeId,
			&barcode,
			&status,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Sales = append(resp.Sales, &models.Sale{
			Id:          id.String,
			SaleID:      saleId.String,
			BranchID:    branchId.String,
			SalePointID: salepointId.String,
			ShiftID:     shiftId.String,
			EmployeeID:  employeeId.String,
			Barcode:     barcode.String,
			Status:      status.String,
			CreatedAt:   createdAt.String,
			UpdatedAt:   updatedAt.String,
		})
	}

	return &resp, nil
}

func (r *saleRepo) Update(ctx context.Context, req *models.UpdateSale) (int64, error) {

	query := `
		UPDATE sale
			SET
				sale_id = $2,
				branch_id = $3,
				salepoint_id = $4,
				shift_id = $5,
				employee_id = $6,
				barcode = $7,
				status = $8,
				updated_at = NOW()
		WHERE id = $1
	`
	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		helpers.NewNullString(req.BranchID),
		helpers.NewNullString(req.SalePointID),
		helpers.NewNullString(req.ShiftID),
		helpers.NewNullString(req.EmployeeID),
		req.Barcode,
		req.Status,
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *saleRepo) Delete(ctx context.Context, req *models.SalePrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM sale WHERE id = $1", req.Id)
	return err
}
