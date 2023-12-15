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

type shiftRepo struct {
	db *pgxpool.Pool
}

func NewShiftRepo(db *pgxpool.Pool) *shiftRepo {
	return &shiftRepo{
		db: db,
	}
}

func (r *shiftRepo) Create(ctx context.Context, req *models.CreateShift) (*models.Shift, error) {
	var (
		shiftId = uuid.New().String()
		query   = `
			INSERT INTO shift(
				id,
				branch_id,
				employee_id,
				sale_point_id,
				status,
				open_shift,
				close_shift,
				updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW());
		`
	)

	_, err := r.db.Exec(ctx,
		query,
		shiftId,
		helpers.NewNullString(req.BranchID),
		helpers.NewNullString(req.EmployeeID),
		helpers.NewNullString(req.SalePointID),
		req.Status,
		helpers.NewNullString(req.OpenShift),
		helpers.NewNullString(req.CloseShift),
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.ShiftPrimaryKey{Id: shiftId})
}

func (r *shiftRepo) GetByID(ctx context.Context, req *models.ShiftPrimaryKey) (*models.Shift, error) {
	var (
		query = `
			SELECT
				id,
				branch_id,
				employee_id,
				sale_point_id,
				status,
				open_shift,
				close_shift,
				created_at,
				updated_at
			FROM shift
			WHERE id = $1
		`
	)

	var (
		Id          sql.NullString
		BranchID    sql.NullString
		EmployeeID  sql.NullString
		SalePointID sql.NullString
		Status      sql.NullString
		OpenShift   sql.NullString
		CloseShift  sql.NullString
		CreatedAt   sql.NullString
		UpdatedAt   sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&Id,
		&BranchID,
		&EmployeeID,
		&SalePointID,
		&Status,
		&OpenShift,
		&CloseShift,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Shift{
		Id:          Id.String,
		BranchID:    BranchID.String,
		EmployeeID:  EmployeeID.String,
		SalePointID: SalePointID.String,
		Status:      Status.String,
		OpenShift:   OpenShift.String,
		CloseShift:  CloseShift.String,
		CreatedAt:   CreatedAt.String,
		UpdatedAt:   UpdatedAt.String,
	}, nil
}

func (r *shiftRepo) GetList(ctx context.Context, req *models.GetListShiftRequest) (*models.GetListShiftResponse, error) {
	var (
		resp   models.GetListShiftResponse
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
		where += " AND (branch_id ILIKE" + " '%" + req.Search + "%' OR employee_id ILIKE '%" + req.Search + "%')"
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			id,
			branch_id,
			employee_id,
			sale_point_id,
			status,
			open_shift,
			close_shift,
			created_at,
			updated_at
		FROM shift
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var (
			Id          sql.NullString
			BranchID    sql.NullString
			EmployeeID  sql.NullString
			SalePointID sql.NullString
			Status      sql.NullString
			OpenShift   sql.NullString
			CloseShift  sql.NullString
			CreatedAt   sql.NullString
			UpdatedAt   sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&Id,
			&BranchID,
			&EmployeeID,
			&SalePointID,
			&Status,
			&OpenShift,
			&CloseShift,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Shift = append(resp.Shift, &models.Shift{
			Id:          Id.String,
			BranchID:    BranchID.String,
			EmployeeID:  EmployeeID.String,
			SalePointID: SalePointID.String,
			Status:      Status.String,
			OpenShift:   OpenShift.String,
			CloseShift:  CloseShift.String,
			CreatedAt:   CreatedAt.String,
			UpdatedAt:   UpdatedAt.String,
		})
	}

	return &resp, nil
}

func (r *shiftRepo) Update(ctx context.Context, req *models.UpdateShift) (int64, error) {
	query := `
		UPDATE shift
			SET
				branch_id = $2,
				employee_id = $3,
				sale_point_id = $4,
				status = $5,
				open_shift = $6,
				close_shift = $7,
				updated_at = NOW()
		WHERE id = $1
	`
	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		helpers.NewNullString(req.BranchID),
		helpers.NewNullString(req.EmployeeID),
		helpers.NewNullString(req.SalePointID),
		helpers.NewNullString(req.Status),
		helpers.NewNullString(req.OpenShift),
		helpers.NewNullString(req.CloseShift),
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *shiftRepo) Delete(ctx context.Context, req *models.ShiftPrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM shift WHERE id = $1", req.Id)
	return err
}
