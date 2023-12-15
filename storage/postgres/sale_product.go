package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"market_system/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type saleProductRepo struct {
	db *pgxpool.Pool
}

func NewSaleProductRepo(db *pgxpool.Pool) *saleProductRepo {
	return &saleProductRepo{
		db: db,
	}
}

func (r *saleProductRepo) Create(ctx context.Context, req *models.CreateSaleProduct) (*models.SaleProduct, error) {
	var (
		saleProductId = uuid.New().String()
		query         = `
			INSERT INTO sale_products(
				id,
				sale_id,
				category_id,
				product_name,
				barcode,
				remaining_quantity,
				quantity,
				allow_discount,
				discount_type,
				discount,
				price,
				total_amount,
				updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW())
		`
	)

	_, err := r.db.Exec(ctx,
		query,
		saleProductId,
		req.SaleID,
		req.CategoryID,
		req.ProductName,
		req.Barcode,
		req.RemainingQuantity,
		req.Quantity,
		req.AllowDiscount,
		req.DiscountType,
		req.Discount,
		req.Price,
		req.TotalAmount,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.SaleProductPrimaryKey{Id: saleProductId})
}

func (r *saleProductRepo) GetByID(ctx context.Context, req *models.SaleProductPrimaryKey) (*models.SaleProduct, error) {
	var (
		query = `
			SELECT
				id,
				sale_id,
				category_id,
				product_name,
				barcode,
				remaining_quantity,
				quantity,
				allow_discount,
				discount_type,
				discount,
				price,
				total_amount,
				created_at,
				updated_at
			FROM sale_products
			WHERE id = $1
		`
	)

	var (
		ID                sql.NullString
		SaleID            sql.NullString
		CategoryID        sql.NullString
		ProductName       sql.NullString
		Barcode           sql.NullString
		RemainingQuantity sql.NullInt64
		Quantity          sql.NullInt64
		AllowDiscount     sql.NullBool
		DiscountType      sql.NullString
		Discount          sql.NullFloat64
		Price             sql.NullFloat64
		TotalAmount       sql.NullFloat64
		CreatedAt         sql.NullString
		UpdatedAt         sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&ID,
		&SaleID,
		&CategoryID,
		&ProductName,
		&Barcode,
		&RemainingQuantity,
		&Quantity,
		&AllowDiscount,
		&DiscountType,
		&Discount,
		&Price,
		&TotalAmount,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.SaleProduct{
		Id:                ID.String,
		SaleID:            SaleID.String,
		CategoryID:        CategoryID.String,
		ProductName:       ProductName.String,
		Barcode:           Barcode.String,
		RemainingQuantity: int(RemainingQuantity.Int64),
		Quantity:          int(Quantity.Int64),
		AllowDiscount:     AllowDiscount.Bool,
		DiscountType:      DiscountType.String,
		Discount:          Discount.Float64,
		Price:             Price.Float64,
		TotalAmount:       TotalAmount.Float64,
		CreatedAt:         CreatedAt.String,
		UpdatedAt:         UpdatedAt.String,
	}, nil
}

func (r *saleProductRepo) GetList(ctx context.Context, req *models.GetListSaleProductRequest) (*models.GetListSaleProductResponse, error) {
	var (
		resp   models.GetListSaleProductResponse
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
			sale_id,
			category_id,
			product_name,
			barcode,
			remaining_quantity,
			quantity,
			allow_discount,
			discount_type,
			discount,
			price,
			total_amount,
			created_at,
			updated_at
		FROM sale_products
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			ID                sql.NullString
			SaleID            sql.NullString
			CategoryID        sql.NullString
			ProductName       sql.NullString
			Barcode           sql.NullString
			RemainingQuantity sql.NullInt64
			Quantity          sql.NullInt64
			AllowDiscount     sql.NullBool
			DiscountType      sql.NullString
			Discount          sql.NullFloat64
			Price             sql.NullFloat64
			TotalAmount       sql.NullFloat64
			CreatedAt         sql.NullString
			UpdatedAt         sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&ID,
			&SaleID,
			&CategoryID,
			&ProductName,
			&Barcode,
			&RemainingQuantity,
			&Quantity,
			&AllowDiscount,
			&DiscountType,
			&Discount,
			&Price,
			&TotalAmount,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		resp.SaleProducts = append(resp.SaleProducts, &models.SaleProduct{
			Id:                ID.String,
			SaleID:            SaleID.String,
			CategoryID:        CategoryID.String,
			ProductName:       ProductName.String,
			Barcode:           Barcode.String,
			RemainingQuantity: int(RemainingQuantity.Int64),
			Quantity:          int(Quantity.Int64),
			AllowDiscount:     AllowDiscount.Bool,
			DiscountType:      DiscountType.String,
			Discount:          Discount.Float64,
			Price:             Price.Float64,
			TotalAmount:       TotalAmount.Float64,
			CreatedAt:         CreatedAt.String,
			UpdatedAt:         UpdatedAt.String,
		})
	}

	return &resp, nil
}

func (r *saleProductRepo) Update(ctx context.Context, req *models.UpdateSaleProduct) (int64, error) {
	query := `
		UPDATE sale_products
		SET
			remaining_quantity = $2,
			quantity = $3,
			allow_discount = $4,
			discount_type = $5,
			discount = $6,
			price = $7,
			total_amount = $8,
			updated_at = NOW()
		WHERE id = $1
	`
	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		req.RemainingQuantity,
		req.Quantity,
		req.AllowDiscount,
		req.DiscountType,
		req.Discount,
		req.Price,
		req.TotalAmount,
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *saleProductRepo) Delete(ctx context.Context, req *models.SaleProductPrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM sale_products WHERE id = $1", req.Id)
	return err
}
