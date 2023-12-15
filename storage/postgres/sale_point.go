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

type SalePointRepo struct {
	db *pgxpool.Pool
}

func NewSalePointRepo(db *pgxpool.Pool) *SalePointRepo {
	return &SalePointRepo{
		db: db,
	}
}

func (r *SalePointRepo) Create(ctx context.Context, req *models.CreateSalePoint) (*models.SalePoint, error) {

	var (
		salePointID = uuid.New().String()
		query       = `
			INSERT INTO sale_point(
				id,
				branch_id,
				name,
				created_at,
				updated_at
			) VALUES ($1, $2, $3, NOW(), NOW())`
	)

	_, err := r.db.Exec(ctx,
		query,
		salePointID,
		helpers.NewNullString(req.Branch_id),
		req.Name,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.SalePointPrimaryKey{Id: salePointID})
}

func (r *SalePointRepo) GetByID(ctx context.Context, req *models.SalePointPrimaryKey) (*models.SalePoint, error) {

	var (
		query = `
			SELECT
				id,
				branch_id,
				name,
				created_at,
				updated_at	
			FROM sale_point
			WHERE id = $1
		`
	)

	var (
		Id        sql.NullString
		BranchID  sql.NullString
		Name      sql.NullString
		CreatedAt sql.NullString
		UpdatedAt sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&Id,
		&BranchID,
		&Name,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.SalePoint{
		Id:        Id.String,
		Branch_id: BranchID.String,
		Name:      Name.String,
		CreatedAt: CreatedAt.String,
		UpdatedAt: UpdatedAt.String,
	}, nil
}

func (r *SalePointRepo) GetList(ctx context.Context, req *models.GetListSalePointRequest) (*models.GetListSalePointResponse, error) {
	var (
		resp   models.GetListSalePointResponse
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
		where += fmt.Sprintf(" AND name ILIKE '%%%s%%'", req.Search)
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			id,
			branch_id,
			name,
			created_at,
			updated_at
		FROM sale_point
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			Id        sql.NullString
			BranchID  sql.NullString
			Name      sql.NullString
			CreatedAt sql.NullString
			UpdatedAt sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&Id,
			&BranchID,
			&Name,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		resp.SalePoints = append(resp.SalePoints, &models.SalePoint{
			Id:        Id.String,
			Branch_id: BranchID.String,
			Name:      Name.String,
			CreatedAt: CreatedAt.String,
			UpdatedAt: UpdatedAt.String,
		})

	}

	return &resp, nil
}

func (r *SalePointRepo) Update(ctx context.Context, req *models.UpdateSalePoint) (int64, error) {

	query := `
		UPDATE sale_point
			SET
				name = $2,
				branch_id = $3,
				updated_at = NOW()
		WHERE id = $1
	`
	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		req.Name,
		helpers.NewNullString(req.Branch_id),
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *SalePointRepo) Delete(ctx context.Context, req *models.SalePointPrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM sale_point WHERE id = $1", req.Id)
	return err
}
