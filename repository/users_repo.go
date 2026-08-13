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
	GetUserQuery := `SELECT id,hashed_password,full_name,created_at,updated_at,role FROM users WHERE email=$1;`
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
	GetUserQuery := `SELECT email,hashed_password,full_name,created_at,updated_at,role FROM users WHERE id=$1;`
	err := UR.DB.DB.QueryRowContext(ctx, GetUserQuery, id).Scan(&user.Email, &user.HashedPassword, &user.FullName, &user.CreatedAt, &user.UpdatedAt, &user.Role)
	if err == sql.ErrNoRows {
		return models.User{}, apperrors.ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}
	user.Id = id
	return user, nil
}

func (UP *UserRepository) UpdateUserFullName(ctx context.Context, id int, NewName string) (models.User, error) {
	var user models.User
	UpdateQuery := `UPDATE users SET full_name = $1 WHERE id=$2 RETURNING email,hashed_password,created_at,updated_at,role;`
	err := UP.DB.DB.QueryRowContext(ctx, UpdateQuery, NewName, id).Scan(&user.Email, &user.HashedPassword, &user.CreatedAt, &user.UpdatedAt, &user.Role)
	if err == sql.ErrNoRows {
		return models.User{}, apperrors.ErrUserNotFound
	} else if err != nil {
		return models.User{}, err
	}
	user.FullName = NewName
	user.Id = id
	return user, nil
}

func (UP *UserRepository) UpdateUserEmail(ctx context.Context, id int, email string) (models.User, error) {
	var user models.User

	UpdateQuery := `UPDATE users SET email = $1 WHERE id=$2 RETURNING full_name,hashed_password,created_at,updated_at,role;`
	err := UP.DB.DB.QueryRowContext(ctx, UpdateQuery, email, id).Scan(&user.FullName, &user.HashedPassword, &user.CreatedAt, &user.UpdatedAt, &user.Role)
	if err == sql.ErrNoRows {
		return models.User{}, apperrors.ErrUserNotFound
	} else if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return models.User{}, apperrors.ErrEmailIsUsed
			}
		}
		return models.User{}, err
	}
	user.Email = email
	user.Id = id
	return user, nil
}

func (UP *UserRepository) UpdateUserHashedPassword(ctx context.Context, NewHashedPassword string, id int) (models.User, error) {
	var user models.User
	UpdateQuery := `UPDATE users SET hashed_password=$1 WHERE id=$2 RETURNING email,updated_at,created_at,full_name,role;`
	err := UP.DB.DB.QueryRowContext(ctx, UpdateQuery, NewHashedPassword, id).Scan(&user.Email, &user.UpdatedAt, &user.CreatedAt, &user.FullName, &user.Role)
	if err == sql.ErrNoRows {
		return models.User{}, apperrors.ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}
	user.Id = id
	user.HashedPassword = NewHashedPassword
	return user, nil
}
