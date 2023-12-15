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

type categoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) *categoryRepo {
	return &categoryRepo{
		db: db,
	}
}

func (r *categoryRepo) Create(ctx context.Context, req *models.CreateCategory) (*models.Category, error) {

	var (
		categoryId = uuid.New().String()
		query      = `
			INSERT INTO category(
				id,
				title, 
				parent_id, 
				updated_at
			) VALUES ($1, $2, $3, NOW())`
	)

	_, err := r.db.Exec(ctx,
		query,
		categoryId,
		req.Title,
		helpers.NewNullString(req.ParentID),
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.CategoryPrimaryKey{Id: categoryId})
}

func (r *categoryRepo) GetByID(ctx context.Context, req *models.CategoryPrimaryKey) (*models.Category, error) {

	var (
		query = `
			SELECT
				id,
				title,
				COALESCE(CAST(parent_id AS VARCHAR), ''),
				created_at,
				updated_at	
			FROM category
			WHERE id = $1
		`
	)

	var (
		Id        sql.NullString
		Title     sql.NullString
		ParentID  sql.NullString
		CreatedAt sql.NullString
		UpdatedAt sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&Id,
		&Title,
		&ParentID,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Category{
		Id:        Id.String,
		Title:     Title.String,
		ParentID:  ParentID.String,
		CreatedAt: CreatedAt.String,
		UpdatedAt: UpdatedAt.String,
	}, nil
}

func (r *categoryRepo) GetList(ctx context.Context, req *models.GetListCategoryRequest) (*models.GetListCategoryResponse, error) {
	var (
		resp   models.GetListCategoryResponse
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
		where += " AND title ILIKE" + " '%" + req.Search + "%'"
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			id,
			title,
			parent_id,
			created_at,
			updated_at
		FROM category
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			Id        sql.NullString
			Title     sql.NullString
			ParentID  sql.NullString
			CreatedAt sql.NullString
			UpdatedAt sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&Id,
			&Title,
			&ParentID,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Categories = append(resp.Categories, &models.Category{
			Id:        Id.String,
			Title:     Title.String,
			ParentID:  ParentID.String,
			CreatedAt: CreatedAt.String,
			UpdatedAt: UpdatedAt.String,
		})
	}

	return &resp, nil
}

func (r *categoryRepo) Update(ctx context.Context, req *models.UpdateCategory) (int64, error) {

	query := `
		UPDATE category
			SET
				title = $2,
				parent_id = $3,
				updated_at = NOW()
		WHERE id = $1
	`
	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		req.Title,
		helpers.NewNullString(req.ParentID),
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *categoryRepo) Delete(ctx context.Context, req *models.CategoryPrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM category WHERE id = $1", req.Id)
	return err
}
