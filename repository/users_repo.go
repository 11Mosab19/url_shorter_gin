package repository

import (
	"context"
	"database/sql"
	"errors"
	"url_shorter_gin/apperrors"
	db "url_shorter_gin/database"
	"url_shorter_gin/models"

	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository struct {
	DB *db.Database
}

func (UR *UserRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	var id int
	CreateUserQuery := `INSERT INTO users (email,hashed_password,full_name) VALUES ($1,$2,$3) RETURNING id;`
	err := UR.DB.DB.QueryRowContext(ctx, CreateUserQuery, user.Email, user.HashedPassword, user.FullName).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return models.User{}, apperrors.ErrEmailIsUsed
			}
		}
		return models.User{}, err
	}
	GetUserData := `SELECT created_at,updated_at,role FROM users WHERE id=$1;`
	err = UR.DB.DB.QueryRowContext(ctx, GetUserData, id).Scan(&user.CreatedAt, &user.UpdatedAt, &user.Role)
	if err == sql.ErrNoRows {
		return models.User{}, apperrors.ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}
	user.Id = id
	return user, nil
}

func (UR *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	GetUserQuery := `SELECT id,hashed_password,full_name,created_at,updated_at,role WHERE email=$1;`
	err := UR.DB.DB.QueryRowContext(ctx, GetUserQuery, email).Scan(&user.Id, &user.HashedPassword, &user.FullName, &user.CreatedAt, &user.UpdatedAt, &user.Role)
	if err == sql.ErrNoRows {
		return models.User{}, apperrors.ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}
	user.Email = email
	return user, nil
}

func (UR *UserRepository) GetUserById(ctx context.Context, id int) (models.User, error) {
	var user models.User
	GetUserQuery := `SELECT email,hashed_password,full_name,created_at,updated_at,role WHERE id=$1;`
	err := UR.DB.DB.QueryRowContext(ctx, GetUserQuery, id).Scan(&user.Id, &user.HashedPassword, &user.FullName, &user.CreatedAt, &user.UpdatedAt, &user.Role)
	if err == sql.ErrNoRows {
		return models.User{}, apperrors.ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}
	user.Id = id
	return user, nil
}
