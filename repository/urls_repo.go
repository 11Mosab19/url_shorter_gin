package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"url_shorter_gin/apperrors"
	db "url_shorter_gin/database"
	"url_shorter_gin/models"

	"github.com/jackc/pgx/v5/pgconn"
)

type UrlRepository struct {
	DB *db.Database
}

func (UrR *UrlRepository) CreateUrl(ctx context.Context, url models.Url) (models.Url, error) {
	InsertQuery := `INSERT INTO urls (original_url,short_code,expires_at,hashed_password,user_id) VALUES ($1,$2,$3,$4,$5) RETURNING id,created_at,updated_at,total_clicks,status;`
	err := UrR.DB.DB.QueryRowContext(ctx, InsertQuery, url.OriginalUrl, url.ShortCode, url.ExpiresAt, url.HashedPassword, url.UserID).Scan(&url.Id, &url.CreatedAt, &url.UpdatedAt, &url.TotalClicks, &url.Status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return models.Url{}, apperrors.ErrShortCodeUsed
			}
		}
		return models.Url{}, err
	}
	return url, nil
}

func (UrR *UrlRepository) GetUrlById(ctx context.Context, id int) (models.Url, error) {
	var url models.Url
	GetUrlQuery := `SELECT original_url,short_code,created_at,expires_at,status,hashed_password,updated_at,total_clicks,user_id FROM urls WHERE id=$1;`
	err := UrR.DB.DB.QueryRowContext(ctx, GetUrlQuery, id).Scan(&url.OriginalUrl, &url.ShortCode, &url.CreatedAt, &url.ExpiresAt, &url.Status, &url.HashedPassword, &url.UpdatedAt, &url.TotalClicks, &url.UserID)
	if err == sql.ErrNoRows {
		return models.Url{}, apperrors.ErrUrlNotFound
	}
	if err != nil {
		return models.Url{}, err
	}
	url.Id = id
	return url, nil
}

func (UrR *UrlRepository) GetUrlByShortCode(ctx context.Context, code string) (models.Url, error) {
	var url models.Url
	GetUrlQuery := `SELECT original_url,id,created_at,expires_at,status,hashed_password,updated_at,total_clicks,user_id FROM urls WHERE short_code=$1;`
	err := UrR.DB.DB.QueryRowContext(ctx, GetUrlQuery, code).Scan(&url.OriginalUrl, &url.Id, &url.CreatedAt, &url.ExpiresAt, &url.Status, &url.HashedPassword, &url.UpdatedAt, &url.TotalClicks, &url.UserID)
	if err == sql.ErrNoRows {
		return models.Url{}, apperrors.ErrUrlNotFound
	}
	if err != nil {
		return models.Url{}, err
	}
	url.ShortCode = code
	return url, nil
}

func (UrR *UrlRepository) GetUserUrl(ctx context.Context, UserId int) ([]models.Url, error) {
	var urls []models.Url
	GetURlsQuery := `SELECT original_url,short_code,id,created_at,expires_at,status,hashed_password,updated_at,total_clicks,user_id FROM urls WHERE user_id = $1 AND status != 'deleted';`
	rows, err := UrR.DB.DB.QueryContext(ctx, GetURlsQuery, UserId)
	if err != nil {
		return []models.Url{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var url models.Url
		err := rows.Scan(&url.OriginalUrl, &url.ShortCode, &url.Id, &url.CreatedAt, &url.ExpiresAt, &url.Status, &url.HashedPassword, &url.UpdatedAt, &url.TotalClicks, &url.UserID)
		if err != nil {
			return []models.Url{}, err
		}
		urls = append(urls, url)
	}
	if err := rows.Err(); err != nil {
		return []models.Url{}, err
	}
	return urls, nil
}

func (UrR *UrlRepository) DeleteUrlById(ctx context.Context, id int, UserId int) error {
	DeleteQuery := `UPDATE urls SET status = 'deleted' WHERE id =$1 AND user_id=$2;`
	result, err := UrR.DB.DB.ExecContext(ctx, DeleteQuery, id, UserId)
	if err != nil {
		return err
	}
	if x, err1 := result.RowsAffected(); err1 != nil {
		return err1
	} else if x == 0 {
		return apperrors.ErrUrlNotFound
	}
	return nil
}

func (UrR *UrlRepository) GetByOriginalURLAndUser(ctx context.Context, UserId int, url string) error {
	var id int
	GetUrlQuery := `SELECT id FROM urls WHERE original_url=$1 AND user_id=$2 AND status!='deleted';`
	result := UrR.DB.DB.QueryRowContext(ctx, GetUrlQuery, url, UserId).Scan(&id)
	if result == sql.ErrNoRows {
		return nil
	}
	if result != nil {
		return result
	}
	return apperrors.ErrAlreadyExistUrl
}

func (UrR *UrlRepository) UpdateUrlStatus(ctx context.Context, status string, id int, UserId int) (models.Url, error) {
	var url models.Url
	UpdateQuery := `UPDATE urls SET status = $1 WHERE id=$2 AND user_id = $3 RETURNING original_url,short_code,created_at,expires_at,hashed_password,total_clicks,user_id;`
	err := UrR.DB.DB.QueryRowContext(ctx, UpdateQuery, status, id, UserId).Scan(&url.OriginalUrl, &url.ShortCode, &url.CreatedAt, &url.ExpiresAt, &url.HashedPassword, &url.TotalClicks, &url.UserID)
	if err == sql.ErrNoRows {
		return models.Url{}, apperrors.ErrUrlNotFound
	} else if err != nil {
		return models.Url{}, err
	}
	url.Status = status
	url.Id = id
	return url, nil
}

func (UrR *UrlRepository) UpdateUrlHashedPassword(ctx context.Context, id int, NewHashedPassword string, UserId int) (models.Url, error) {
	var url models.Url
	UpdateQuery := `UPDATE urls SET  hashed_password= $1 WHERE id=$2 AND user_id = $3 RETURNING original_url,short_code,created_at,expires_at,status,total_clicks,user_id;`
	err := UrR.DB.DB.QueryRowContext(ctx, UpdateQuery, NewHashedPassword, id, UserId).Scan(&url.OriginalUrl, &url.ShortCode, &url.CreatedAt, &url.ExpiresAt, &url.Status, &url.TotalClicks, &url.UserID)
	if err == sql.ErrNoRows {
		return models.Url{}, apperrors.ErrUrlNotFound
	}
	if err != nil {
		return models.Url{}, err
	}
	url.HashedPassword = NewHashedPassword
	url.Id = id
	return url, nil
}

func (UrR *UrlRepository) UpdateUrlExpireDate(ctx context.Context, id int, UserId int, NewDate *time.Time) (models.Url, error) {
	var url models.Url
	UpdateQuery := `UPDATE urls SET expires_at = $1 WHERE id=$2 AND user_id = $3 RETURNING original_url,short_code,created_at,status,hashed_password,total_clicks,user_id;`
	err := UrR.DB.DB.QueryRowContext(ctx, UpdateQuery, NewDate, id, UserId).Scan(&url.OriginalUrl, &url.ShortCode, &url.CreatedAt, &url.Status, &url.HashedPassword, &url.TotalClicks, &url.UserID)
	if err == sql.ErrNoRows {
		return models.Url{}, apperrors.ErrUrlNotFound
	} else if err != nil {
		return models.Url{}, err
	}
	url.ExpiresAt = NewDate
	url.Id = id
	return url, nil
}
