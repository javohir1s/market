package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"market_system/config"
	"market_system/models"
	"market_system/pkg/helpers"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, req *models.CreateUser) (*models.User, error) {

	if !helpers.Contains(config.ClientTypes, req.ClientType) {
		return nil, errors.New("not found client type")
	}

	var (
		userId = uuid.New().String()
		query  = `
			INSERT INTO "user"(
				"id",
				"first_name",
				"last_name",
				"login",
				"password",
				"active",
				"client_type",
				"updated_at"
			) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`
	)

	_, err := r.db.Exec(ctx,
		query,
		userId,
		req.FirstName,
		req.LastName,
		req.Login,
		req.Password,
		req.Active,
		req.ClientType,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.UserPrimaryKey{Id: userId})
}

func (r *userRepo) GetByID(ctx context.Context, req *models.UserPrimaryKey) (*models.User, error) {

	var (
		query = `
			SELECT
				id,
				first_name,
				last_name,
				login,
				password,
				active,
				client_type,
				created_at,
				updated_at	
			FROM "user"
		`
		where = "WHERE id = $1"
	)

	if !helpers.IsValidUUID(req.Id) {
		where = "WHERE login = $1"
	}

	var (
		Id         sql.NullString
		FirstName  sql.NullString
		LastName   sql.NullString
		Login      sql.NullString
		Password   sql.NullString
		Active     sql.NullBool
		ClientType sql.NullString
		CreatedAt  sql.NullString
		UpdatedAt  sql.NullString
	)

	query += where
	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&Id,
		&FirstName,
		&LastName,
		&Login,
		&Password,
		&Active,
		&ClientType,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.User{
		Id:         Id.String,
		FirstName:  FirstName.String,
		LastName:   LastName.String,
		Login:      Login.String,
		Password:   Password.String,
		Active:     Active.Bool,
		ClientType: ClientType.String,
		CreatedAt:  CreatedAt.String,
		UpdatedAt:  UpdatedAt.String,
	}, nil
}

func (r *userRepo) GetList(ctx context.Context, req *models.GetListUserRequest) (*models.GetListUserResponse, error) {
	var (
		resp   models.GetListUserResponse
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

	// if len(req.Search) > 0 {
	// 	where += " AND first_name ILIKE" + " '%" + req.Search + "%'"
	// }

	var query = `
		SELECT
			COUNT(*) OVER(),
			id,
			first_name,
			last_name,
			login,
			password,
			active,
			client_type,
			created_at,
			updated_at	
		FROM "user"
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			Id         sql.NullString
			FirstName  sql.NullString
			LastName   sql.NullString
			Login      sql.NullString
			Password   sql.NullString
			Active     sql.NullBool
			ClientType sql.NullString
			CreatedAt  sql.NullString
			UpdatedAt  sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&Id,
			&FirstName,
			&LastName,
			&Login,
			&Password,
			&Active,
			&ClientType,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.User = append(resp.User, &models.User{
			Id:         Id.String,
			FirstName:  FirstName.String,
			LastName:   LastName.String,
			Login:      Login.String,
			Password:   Password.String,
			Active:     Active.Bool,
			ClientType: ClientType.String,
			CreatedAt:  CreatedAt.String,
			UpdatedAt:  UpdatedAt.String,
		})

	}

	return &resp, nil
}

func (r *userRepo) Update(ctx context.Context, req *models.UpdateUser) (int64, error) {

	query := `
		UPDATE user
			SET
				first_name = $2,
				last_name = $3,
				login = $4,
				password = $5,
				active = $6,
				client_type = $7,
				updated_at = NOW()
		WHERE id = $1
	`
	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		req.FirstName,
		req.LastName,
		req.Login,
		req.Password,
		req.Active,
		req.ClientType,
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *userRepo) Delete(ctx context.Context, req *models.UserPrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM user WHERE id = $1", req.Id)
	return err
}
