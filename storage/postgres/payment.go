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

type paymentRepo struct {
	db *pgxpool.Pool
}

func NewPaymentRepo(db *pgxpool.Pool) *paymentRepo {
	return &paymentRepo{
		db: db,
	}
}

func (r *paymentRepo) Create(ctx context.Context, req *models.CreatePayment) (*models.Payment, error) {
	paymentId := uuid.New().String()
	query := `
		INSERT INTO payment (
			id,
			sale_id,
			cash,
			uzcard,
			payme,
			click,
			humo,
			apelsin,
			total_amount,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
	`

	_, err := r.db.Exec(ctx,
		query,
		paymentId,
		helpers.NewNullString(req.SaleID),
		req.Cash,
		req.Uzcard,
		req.Payme,
		req.Click,
		req.Humo,
		req.Apelsin,
		req.TotalAmount,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.PaymentPrimaryKey{Id: paymentId})
}

func (r *paymentRepo) GetByID(ctx context.Context, req *models.PaymentPrimaryKey) (*models.Payment, error) {
	query := `
		SELECT
			id,
			sale_id,
			cash,
			uzcard,
			payme,
			click,
			humo,
			apelsin,
			total_amount,
			created_at,
			updated_at
		FROM payment
		WHERE id = $1
	`
	var (
		Id          sql.NullString
		SaleID      sql.NullString
		Cash        sql.NullFloat64
		Uzcard      sql.NullFloat64
		Payme       sql.NullFloat64
		Click       sql.NullFloat64
		Humo        sql.NullFloat64
		Apelsin     sql.NullFloat64
		TotalAmount sql.NullFloat64
		CreatedAt   sql.NullString
		UpdatedAt   sql.NullString
	)
	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&Id,
		&SaleID,
		&Cash,
		&Uzcard,
		&Payme,
		&Click,
		&Humo,
		&Apelsin,
		&TotalAmount,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Payment{
		Id:          Id.String,
		SaleID:      SaleID.String,
		Cash:        Cash.Float64,
		Uzcard:      Uzcard.Float64,
		Payme:       Payme.Float64,
		Click:       Click.Float64,
		Humo:        Humo.Float64,
		Apelsin:     Apelsin.Float64,
		TotalAmount: TotalAmount.Float64,
		CreatedAt:   CreatedAt.String,
		UpdatedAt:   UpdatedAt.String,
	}, nil
}

func (r *paymentRepo) GetList(ctx context.Context, req *models.GetListPaymentRequest) (*models.GetListPaymentResponse, error) {
	var (
		resp   models.GetListPaymentResponse
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
		where += fmt.Sprintf(" AND sale_id ILIKE '%%%s%%'", req.Search)
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	query := `
		SELECT
			COUNT(*) OVER(),
			id,
			sale_id,
			cash,
			uzcard,
			payme,
			click,
			humo,
			apelsin,
			total_amount,
			created_at,
			updated_at
		FROM payment
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var (
			Id          sql.NullString
			SaleID      sql.NullString
			Cash        sql.NullFloat64
			Uzcard      sql.NullFloat64
			Payme       sql.NullFloat64
			Click       sql.NullFloat64
			Humo        sql.NullFloat64
			Apelsin     sql.NullFloat64
			TotalAmount sql.NullFloat64
			CreatedAt   sql.NullString
			UpdatedAt   sql.NullString
		)
		err = rows.Scan(
			&resp.Count,
			&Id,
			&SaleID,
			&Cash,
			&Uzcard,
			&Payme,
			&Click,
			&Humo,
			&Apelsin,
			&TotalAmount,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		resp.Payments = append(resp.Payments, &models.Payment{
			Id:          Id.String,
			SaleID:      SaleID.String,
			Cash:        Cash.Float64,
			Uzcard:      Uzcard.Float64,
			Payme:       Payme.Float64,
			Click:       Click.Float64,
			Humo:        Humo.Float64,
			Apelsin:     Apelsin.Float64,
			TotalAmount: TotalAmount.Float64,
			CreatedAt:   CreatedAt.String,
			UpdatedAt:   UpdatedAt.String,
		})
	}

	return &resp, nil
}

func (r *paymentRepo) Update(ctx context.Context, req *models.UpdatePayment) (int64, error) {
	query := `
		UPDATE payment
		SET
			cash = $2,
			uzcard = $3,
			payme = $4,
			click = $5,
			humo = $6,
			apelsin = $7,
			total_amount = $8,
			updated_at = NOW()
		WHERE id = $1
	`

	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		req.Cash,
		req.Uzcard,
		req.Payme,
		req.Click,
		req.Humo,
		req.Apelsin,
		req.TotalAmount,
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *paymentRepo) Delete(ctx context.Context, req *models.PaymentPrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM payment WHERE id = $1", req.Id)
	return err
}
