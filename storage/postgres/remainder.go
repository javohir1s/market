package postgres

import (
	"context"
	"database/sql"
	"market_system/models"
	"market_system/pkg/helpers"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type remainderRepo struct {
	db *pgxpool.Pool
}

func NewRemainderRepo(db *pgxpool.Pool) *remainderRepo {
	return &remainderRepo{
		db: db,
	}
}

func (r *remainderRepo) Create(ctx context.Context, req *models.CreateRemainder) (*models.Remainder, error) {

	var (
		remainderID = uuid.New().String()
		query       = `
			INSERT INTO remainder(
				id,
				branch_id,
				category_id,
				product_name,
				barcode,
				price_income,
				quantity,
				updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`
	)

	_, err := r.db.Exec(ctx,
		query,
		remainderID,
		helpers.NewNullString(req.BranchID),
		helpers.NewNullString(req.CategoryID),
		req.ProductName,
		req.Barcode,
		req.PriceIncome,
		req.Quantity,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.RemainderPrimaryKey{Id: remainderID})
}

func (r *remainderRepo) GetByID(ctx context.Context, req *models.RemainderPrimaryKey) (*models.Remainder, error) {

	var (
		query = `
			SELECT
				id,
				branch_id,
				category_id,
				product_name,
				barcode,
				price_income,
				quantity,
				created_at,
				updated_at
			FROM  remainder
			WHERE id = $1
		`
	)

	var (
		ID          sql.NullString
		BranchID    sql.NullString
		CategoryID  sql.NullString
		ProductName sql.NullString
		Barcode     sql.NullString
		PriceIncome sql.NullFloat64
		Quantity    sql.NullInt64
		CreatedAt   sql.NullString
		UpdatedAt   sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&ID,
		&BranchID,
		&CategoryID,
		&ProductName,
		&Barcode,
		&PriceIncome,
		&Quantity,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Remainder{
		Id:          ID.String,
		BranchID:    BranchID.String,
		CategoryID:  CategoryID.String,
		ProductName: ProductName.String,
		Barcode:     Barcode.String,
		PriceIncome: PriceIncome.Float64,
		Quantity:    int(Quantity.Int64),
		CreatedAt:   CreatedAt.String,
		UpdatedAt:   UpdatedAt.String,
	}, nil
}

func (r *remainderRepo) GetList(ctx context.Context, req *models.GetListRemainderRequest) (*models.GetListRemainderResponse, error) {
	var (
		resp  models.GetListRemainderResponse
		where = " WHERE TRUE"
	)

	if len(req.Search) > 0 {
		where += " AND (branch_id ILIKE '%" + req.Search + "%' OR category_id ILIKE '%" + req.Search + "%' OR barcode ILIKE '%" + req.Search + "%')"
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			id,
			branch_id,
			category_id,
			product_name,
			barcode,
			price_income,
			quantity,
			created_at,
			updated_at
		FROM remainder
	`

	query += where
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			ID          sql.NullString
			BranchID    sql.NullString
			CategoryID  sql.NullString
			ProductName sql.NullString
			Barcode     sql.NullString
			PriceIncome sql.NullFloat64
			Quantity    sql.NullInt64
			CreatedAt   sql.NullString
			UpdatedAt   sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&ID,
			&BranchID,
			&CategoryID,
			&ProductName,
			&Barcode,
			&PriceIncome,
			&Quantity,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Remainder = append(resp.Remainder, &models.Remainder{
			Id:          ID.String,
			BranchID:    BranchID.String,
			CategoryID:  CategoryID.String,
			ProductName: ProductName.String,
			Barcode:     Barcode.String,
			PriceIncome: PriceIncome.Float64,
			Quantity:    int(Quantity.Int64),
			CreatedAt:   CreatedAt.String,
			UpdatedAt:   UpdatedAt.String,
		})
	}

	return &resp, nil
}

func (r *remainderRepo) Update(ctx context.Context, req *models.UpdateRemainder) (int64, error) {

	query := `
		UPDATE remainder
			SET
				product_name = $2,
				barcode = $3,
				price_income = $4,
				quantity = $5,
				updated_at = NOW()
		WHERE id = $1
	`
	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		req.ProductName,
		req.Barcode,
		req.PriceIncome,
		req.Quantity,
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *remainderRepo) Delete(ctx context.Context, req *models.RemainderPrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM remainder WHERE id = $1", req.Id)
	return err
}
