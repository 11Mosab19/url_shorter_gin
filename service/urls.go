package service

import (
	"context"
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
	return UrS.GetUserUrls(ctx, UserId)
}

func (UrS *UrlService) UpdateUrlStatus(ctx context.Context, UserId int, data dto.UpdateURLRequest, id int) (models.Url, error) {

}
