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

type transactionRepo struct {
	db *pgxpool.Pool
}

func NewTransactionRepo(db *pgxpool.Pool) *transactionRepo {
	return &transactionRepo{
		db: db,
	}
}

func (r *transactionRepo) Create(ctx context.Context, req *models.CreateTransaction) (*models.Transaction, error) {

	var (
		transactionID = uuid.New().String()
		query         = `
			INSERT INTO transaction(
				id,
				shift_id,
				cash,
				uzcard,
				payme,
				click,
				humo,
				apelsin,
				total_amount,
				created_at,
				updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())`
	)

	_, err := r.db.Exec(ctx,
		query,
		transactionID,
		helpers.NewNullString(req.ShiftID),
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

	return r.GetByID(ctx, &models.TransactionPrimaryKey{Id: transactionID})
}

func (r *transactionRepo) GetByID(ctx context.Context, req *models.TransactionPrimaryKey) (*models.Transaction, error) {

	var (
		query = `
			SELECT
				id,
				shift_id,
				cash,
				uzcard,
				payme,
				click,
				humo,
				apelsin,
				total_amount,
				created_at,
				updated_at
			FROM transaction
			WHERE id = $1
		`
	)

	var (
		ID          sql.NullString
		ShiftID     sql.NullString
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
		&ID,
		&ShiftID,
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

	return &models.Transaction{
		Id:          ID.String,
		ShiftID:     ShiftID.String,
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

func (r *transactionRepo) GetList(ctx context.Context, req *models.GetListTransactonRequest) (*models.GetListTransactionResponse, error) {
	var (
		resp   models.GetListTransactionResponse
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
		where += " AND (shift_id ILIKE '%" + req.Search + "%')"
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			id,
			shift_id,
			cash,
			uzcard,
			payme,
			click,
			humo,
			apelsin,
			total_amount,
			created_at,
			updated_at
		FROM transaction
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			Id          sql.NullString
			shiftID     sql.NullString
			cash        sql.NullFloat64
			uzcard      sql.NullFloat64
			payme       sql.NullFloat64
			click       sql.NullFloat64
			humo        sql.NullFloat64
			apelsin     sql.NullFloat64
			totalAmount sql.NullFloat64
			createdAt   sql.NullString
			updatedAt   sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&Id,
			&shiftID,
			&cash,
			&uzcard,
			&payme,
			&click,
			&humo,
			&apelsin,
			&totalAmount,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Transactions = append(resp.Transactions, &models.Transaction{
			Id:          Id.String,
			ShiftID:     shiftID.String,
			Cash:        cash.Float64,
			Uzcard:      uzcard.Float64,
			Payme:       payme.Float64,
			Click:       click.Float64,
			Humo:        humo.Float64,
			Apelsin:     apelsin.Float64,
			TotalAmount: totalAmount.Float64,
			CreatedAt:   createdAt.String,
			UpdatedAt:   updatedAt.String,
		})
	}

	return &resp, nil
}

func (r *transactionRepo) Update(ctx context.Context, req *models.UpdateTransaction) (int64, error) {

	query := `
		UPDATE transaction
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

func (r *transactionRepo) Delete(ctx context.Context, req *models.TransactionPrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM transaction WHERE id = $1", req.Id)
	return err
}
