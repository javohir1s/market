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

type productRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *productRepo {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) Create(ctx context.Context, req *models.CreateProduct) (*models.Product, error) {

	var (
		productID = uuid.New().String()
		query     = `
			INSERT INTO product(
				id,
				photo,
				title, 
				category_id, 
				barcode,
				price,
				updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, NOW())`
	)

	_, err := r.db.Exec(ctx,
		query,
		productID,
		req.Photo,
		req.Title,
		helpers.NewNullString(req.CategoryID),
		helpers.NewNullString(req.Barcode),
		req.Price,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.ProductPrimaryKey{Id: productID})
}

func (r *productRepo) GetByID(ctx context.Context, req *models.ProductPrimaryKey) (*models.Product, error) {

	var (
		query = `
			SELECT
				id,
				photo,
				title,
				category_id,
				COALESCE(barcode, ''),
				COALESCE(price, 0),
				created_at,
				updated_at
			FROM  product
			WHERE id = $1
		`
	)

	var (
		ID         sql.NullString
		Photo      sql.NullString
		Title      sql.NullString
		CategoryID sql.NullString
		Barcode    sql.NullString
		Price      sql.NullFloat64
		CreatedAt  sql.NullString
		UpdatedAt  sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&ID,
		&Photo,
		&Title,
		&CategoryID,
		&Barcode,
		&Price,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Product{
		Id:         ID.String,
		Photo:      Photo.String,
		Title:      Title.String,
		CategoryID: CategoryID.String,
		Barcode:    Barcode.String,
		Price:      Price.Float64,
		CreatedAt:  CreatedAt.String,
		UpdatedAt:  UpdatedAt.String,
	}, nil
}

func (r *productRepo) GetList(ctx context.Context, req *models.GetListProductRequest) (*models.GetListProductResponse, error) {
	var (
		resp   models.GetListProductResponse
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
		where += " AND (title ILIKE '%" + req.Search + "%' OR barcode ILIKE '%" + req.Search + "%')"
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			id,
			photo,
			title,
			category_id,
			barcode,
			price,
			created_at,
			updated_at
		FROM product
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var (
			ID         sql.NullString
			Photo      sql.NullString
			Title      sql.NullString
			CategoryID sql.NullString
			Barcode    sql.NullString
			Price      sql.NullFloat64
			CreatedAt  sql.NullString
			UpdatedAt  sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&ID,
			&Photo,
			&Title,
			&CategoryID,
			&Barcode,
			&Price,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Products = append(resp.Products, &models.Product{
			Id:         ID.String,
			Photo:      Photo.String,
			Title:      Title.String,
			CategoryID: CategoryID.String,
			Barcode:    Barcode.String,
			Price:      Price.Float64,
			CreatedAt:  CreatedAt.String,
			UpdatedAt:  UpdatedAt.String,
		})
	}

	return &resp, nil
}

func (r *productRepo) Update(ctx context.Context, req *models.UpdateProduct) (int64, error) {

	query := `
		UPDATE product
			SET
				photo = $2,
				title = $3,
				category_id = $4,
				barcode = $5,
				price = $6, 
				updated_at = NOW()
		WHERE id = $1
	`
	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		req.Photo,
		req.Title,
		helpers.NewNullString(req.CategoryID),
		helpers.NewNullString(req.Barcode),
		req.Price,
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *productRepo) Delete(ctx context.Context, req *models.ProductPrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM product WHERE id = $1", req.Id)
	return err
}
