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

type incomeProductRepo struct {
	db *pgxpool.Pool
}

func NewIncomeProductRepo(db *pgxpool.Pool) *incomeProductRepo {
	return &incomeProductRepo{
		db: db,
	}
}

func (r *incomeProductRepo) Create(ctx context.Context, req *models.CreateIncomeProduct) (*models.IncomeProduct, error) {

	var (
		incomeProductId = uuid.New().String()
		query           = `
			INSERT INTO income_product(
				id,
				income_id,
				category_id,
				product_name,
				barcode,
				quantity,
				income_price,
				updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`
	)

	_, err := r.db.Exec(ctx,
		query,
		incomeProductId,
		helpers.NewNullString(req.IncomeID),
		helpers.NewNullString(req.CategoryID),
		req.ProductName,
		req.Barcode,
		req.Quantity,
		req.IncomePrice,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.IncomeProductPrimaryKey{Id: incomeProductId})
}

func (r *incomeProductRepo) GetByID(ctx context.Context, req *models.IncomeProductPrimaryKey) (*models.IncomeProduct, error) {

	var (
		query = `
			SELECT
				id,
				income_id,
				category_id,
				product_name,
				barcode,
				quantity,
				income_price,
				created_at,
				updated_at	
			FROM  income_product
			WHERE id = $1
		`
	)

	var (
		Id          sql.NullString
		IncomeID    sql.NullString
		CategoryID  sql.NullString
		ProductName sql.NullString
		Barcode     sql.NullString
		Quantity    sql.NullInt64
		IncomePrice sql.NullFloat64
		CreatedAt   sql.NullString
		UpdatedAt   sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&Id,
		&IncomeID,
		&CategoryID,
		&ProductName,
		&Barcode,
		&Quantity,
		&IncomePrice,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.IncomeProduct{
		Id:          Id.String,
		IncomeID:    IncomeID.String,
		CategoryID:  CategoryID.String,
		ProductName: ProductName.String,
		Barcode:     Barcode.String,
		Quantity:    Quantity.Int64,
		IncomePrice: IncomePrice.Float64,
		CreatedAt:   CreatedAt.String,
		UpdatedAt:   UpdatedAt.String,
	}, nil
}

func (r *incomeProductRepo) GetList(ctx context.Context, req *models.GetListIncomeProductRequest) (*models.GetListIncomeProductResponse, error) {
	var (
		resp   models.GetListIncomeProductResponse
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
		where += " AND (category_id ILIKE '%" + req.Search + "%' OR barcode ILIKE '%" + req.Search + "%')"
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			 id,
			 income_id,
			 category_id,
			 product_name,
			 barcode,
			 quantity,
			 income_price,
			 created_at,
			 updated_at
		FROM income_product
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			Id          sql.NullString
			IncomeID    sql.NullString
			CategoryID  sql.NullString
			ProductName sql.NullString
			Barcode     sql.NullString
			Quantity    sql.NullInt64
			IncomePrice sql.NullFloat64
			CreatedAt   sql.NullString
			UpdatedAt   sql.NullString
		)

		err = rows.Scan(
			&Id,
			&IncomeID,
			&CategoryID,
			&ProductName,
			&Barcode,
			&Quantity,
			&IncomePrice,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.IncomeProducts = append(resp.IncomeProducts, &models.IncomeProduct{
			Id:          Id.String,
			IncomeID:    IncomeID.String,
			CategoryID:  CategoryID.String,
			ProductName: ProductName.String,
			Barcode:     Barcode.String,
			Quantity:    Quantity.Int64,
			IncomePrice: IncomePrice.Float64,
			CreatedAt:   CreatedAt.String,
			UpdatedAt:   UpdatedAt.String,
		})
	}

	return &resp, nil
}

func (r *incomeProductRepo) Update(ctx context.Context, req *models.UpdateIncomeProduct) (int64, error) {

	query := `
		UPDATE income_product
			SET
				product_name = $2,
				barcode = $3,
				quantity = $4,
				income_price = $5,
				category_id = $6,
				income_id = $7,
				updated_at = NOW()
		WHERE id = $1
	`
	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		req.ProductName,
		req.Barcode,
		req.Quantity,
		req.IncomePrice,
		helpers.NewNullString(req.CategoryID),
		helpers.NewNullString(req.IncomeID),
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *incomeProductRepo) Delete(ctx context.Context, req *models.IncomeProductPrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM income_product WHERE id = $1", req.Id)
	return err
}
