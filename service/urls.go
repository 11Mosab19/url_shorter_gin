package service

import (
	"context"
	"time"
	"url_shorter_gin/apperrors"
	"url_shorter_gin/dto"
	"url_shorter_gin/functions"
	"url_shorter_gin/models"
	"url_shorter_gin/repository"

	"golang.org/x/crypto/bcrypt"
)

type UrlService struct {
	repo repository.UrlRepository
}

func (UrS *UrlService) CreateUrl(ctx context.Context, UrlData dto.CreateURLRequest, UserId int) (models.Url, error) {
	var url models.Url
	if UrlData.ExpiresAt != nil {
		if UrlData.ExpiresAt.Before(time.Now()) {
			return models.Url{}, apperrors.ErrInvalidExpirationInput
		}
	}
	if UrlData.ShortCode == "" {
		GeneratedShortCode, err := functions.GenerateShortCode()
		if err != nil {
			return models.Url{}, apperrors.ErrGenerateShortCode
		}
		url.ShortCode = GeneratedShortCode
	} else {
		url.ShortCode = UrlData.ShortCode
	}
	if UrlData.Password != "" {
		HashedPassword, err := bcrypt.GenerateFromPassword([]byte(UrlData.Password), bcrypt.DefaultCost)
		if err != nil {
			return models.Url{}, apperrors.HashingErr
		}
		url.HashedPassword = string(HashedPassword)
	}
	url.UserID = UserId
	url.OriginalUrl = UrlData.OriginalUrl
	url.ExpiresAt = UrlData.ExpiresAt

	url, err := UrS.repo.CreateUrl(ctx, url)
	if err != nil {
		return models.Url{}, err
	}
	return url, nil
}

func (UrS *UrlService) GetUserUrls(ctx context.Context, UserId int) ([]models.Url, error) {
	return UrS.repo.GetUserUrl(ctx, UserId)
}

func (UrS *UrlService) UpdateUrlStatus(ctx context.Context, UserId int, data dto.UpdateURLRequest, id int) (models.Url, error) {
	url, err := UrS.repo.GetUrlById(ctx, id)
	if err != nil {
		return models.Url{}, err
	}
	if url.UserID != UserId {
		return models.Url{}, apperrors.ErrUnauthorized
	}
	if data.Status == url.Status {
		return url, nil
	}
	if url.Status == "deleted" {
		if data.Status == "disabled" {
			return models.Url{}, apperrors.ErrInvalidStatusInput
		}
		if data.Status == "active" {
			_, err := UrS.repo.UpdateUrlStatus(ctx, data.Status, id, UserId)
			if err != nil {
				return models.Url{}, err
			}
			url.Status = data.Status
			return url, nil
		}
	}
	if url.Status == "active" {
		_, err := UrS.repo.UpdateUrlStatus(ctx, data.Status, id, UserId)
		if err != nil {
			return models.Url{}, err
		}
		url.Status = data.Status
		return url, nil
	}
	if url.Status == "disabled" {
		_, err := UrS.repo.UpdateUrlStatus(ctx, data.Status, id, UserId)
		if err != nil {
			return models.Url{}, err
		}
		url.Status = data.Status
	}
	return url, nil
}

func (Urs *UrlService) RedirectByShortCode(ctx context.Context, ShortCode string, password string) (models.Url, error) {
	url, err := Urs.repo.GetUrlByShortCode(ctx, ShortCode)
	if err != nil {
		return models.Url{}, err
	}
	if url.Status != "active" {
		return models.Url{}, apperrors.ErrBrokenUrl
	}
	if url.ExpiresAt != nil {
		timeout := *url.ExpiresAt
		now := time.Now()
		if timeout.Before(now) {
			return models.Url{}, apperrors.ErrExpiredUrl
		}
	}
	if url.HashedPassword != "" {
		err := bcrypt.CompareHashAndPassword([]byte(url.HashedPassword), []byte(password))
		if err != nil {
			return models.Url{}, apperrors.ErrWrongUrlPassword
		}
	}
	return url, nil
}

func (UrS *UrlService) UpdateUrlExpiration(ctx context.Context, data dto.UpdateURLRequest, id int, userId int) (models.Url, error) {
	if data.Expiration != nil {
		toUpdate := *data.Expiration
		if toUpdate.Before(time.Now()) {
			return models.Url{}, apperrors.ErrInvalidExpirationInput
		}
		url, err := UrS.repo.UpdateUrlExpireDate(ctx, id, userId, data.Expiration)
		if err != nil {
			return models.Url{}, err
		}
		return url, nil
	}
	url, err := UrS.repo.UpdateUrlExpireDate(ctx, id, userId, data.Expiration)
	if err != nil {
		return models.Url{}, err
	}
	return url, nil
}

func (UrS *UrlService) UpdateUrlPassword(ctx context.Context, data dto.UpdateURLRequest, id int, userId int) (models.Url, error) {
	url, err := UrS.repo.GetUrlById(ctx, id)
	if err != nil {
		return models.Url{}, err
	}
	if url.UserID != userId {
		return models.Url{}, apperrors.ErrUnauthorized
	}
	match := bcrypt.CompareHashAndPassword([]byte(url.HashedPassword), []byte(data.OldPassword))
	if match != nil {
		return models.Url{}, apperrors.ErrWrongUrlPassword
	}
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return models.Url{}, apperrors.HashingErr
	}
	url, err = UrS.repo.UpdateUrlHashedPassword(ctx, id, string(HashedPassword), userId)
	if err != nil {
		return models.Url{}, err
	}
	return url, nil
}

func (UrS *UrlService) SetPasswordToUrl(ctx context.Context, Data dto.SetPasswordRequest, id int, userId int) (models.Url, error) {
	if Data.NewPassword != Data.ConfirmationPassword {
		return models.Url{}, apperrors.ErrWrongConfirmationPassword
	}
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(Data.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return models.Url{}, apperrors.HashingErr
	}
	url, err := UrS.repo.UpdateUrlHashedPassword(ctx, id, string(HashedPassword), userId)
	if err != nil {
		return models.Url{}, err
	}
	return url, nil
}

func (UrS *UrlService) DeleteUrlById(ctx context.Context, id int, userId int) error {
	return UrS.repo.DeleteUrlById(ctx, id, userId)
}
