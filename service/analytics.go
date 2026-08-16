package service

import (
	"context"
	"time"
	"url_shorter_gin/models"
	"url_shorter_gin/repository"

	"github.com/medama-io/go-useragent"
)

type AnalyticsService struct {
	repo repository.AnalyticsRepository
}

func (AS *AnalyticsService) CreateAnalytic(ctx context.Context, data models.Analytics) (models.Analytics, error) {
	ua := data.UserAgent
	agent := useragent.NewParser().Parse(ua)

	data.DeviceType = string(agent.Device())
	data.ClickedAt = time.Now()

	analytic, err := AS.repo.CreateAnalytic(ctx, data)
	if err != nil {
		return models.Analytics{}, err
	}
	analytic.UserAgent = data.UserAgent

	return analytic, nil
}

func (AS *AnalyticsService) GetUrlAnalytics(ctx context.Context, urlId int) ([]models.Analytics, error) {
	analytics, err := AS.repo.GetUrlAnalyticsByUrlId(ctx, urlId)
	if err != nil {
		return []models.Analytics{}, err
	}
	return analytics, nil
}

func (AS *AnalyticsService) GetUrlUsersIps(ctx context.Context, urlId int) ([]string, error) {
	Ips, err := AS.repo.GetIpAddressesForUrlById(ctx, urlId)
	if err != nil {
		return []string{}, err
	}
	return Ips, nil
}

func (AS *AnalyticsService) GetUrlUsersDevices(ctx context.Context, urlId int) ([]string, error) {
	devices, err := AS.repo.GetUsedDevicesForUrlWithId(ctx, urlId)
	if err != nil {
		return []string{}, err
	}
	return devices, nil
}
