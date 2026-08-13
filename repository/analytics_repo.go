package repository

import (
	"context"
	"url_shorter_gin/database"
	"url_shorter_gin/models"
)

type AnalyticsRepository struct {
	DB *database.Database
}

func (AnR *AnalyticsRepository) CreateAnalytic(ctx context.Context, analytic models.Analytics) (models.Analytics, error) {
	CreateAnalyticQuery := `INSERT INTO analytics (url_id,clicked_at,device_type,country) VALUES ($1,$2,$3,$4) RETURNING id;`
	err := AnR.DB.DB.QueryRowContext(ctx, CreateAnalyticQuery, analytic.UrlID, analytic.ClickedAt, analytic.DeviceType, analytic.Country).Scan(&analytic.Id)
	if err != nil {
		return models.Analytics{}, err
	}
	return analytic, nil
}

func (AnR *AnalyticsRepository) GetUrlAnalyticsByUrlId(ctx context.Context, UrlId int) ([]models.Analytics, error) {
	var analytics []models.Analytics
	GetAnalyticsQuery := `SELECT id,url_id,clicked_at,device_type,country FROM analytics WHERE url_id = $1 ORDER BY clicked_at DESC; `
	rows, err := AnR.DB.DB.QueryContext(ctx, GetAnalyticsQuery, UrlId)
	if err != nil {
		return []models.Analytics{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var analytic models.Analytics
		err := rows.Scan(&analytic.Id, &analytic.UrlID, &analytic.ClickedAt, &analytic.DeviceType, &analytic.Country)
		if err != nil {
			return []models.Analytics{}, err
		}
		analytics = append(analytics, analytic)
	}
	if err := rows.Err(); err != nil {
		return []models.Analytics{}, err
	}
	return analytics, nil
}

func (AnR *AnalyticsRepository) GetUsedDevicesForUrlWithId(ctx context.Context, UrlId int) ([]string, error) {
	var devices []string
	GetDevicesQuery := `SELECT device_type FROM analytics WHERE url_id =$1;`
	rows, err := AnR.DB.DB.QueryContext(ctx, GetDevicesQuery, UrlId)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var device string
		err := rows.Scan(&device)
		if err != nil {
			return []string{}, err
		}
		devices = append(devices, device)
	}
	if err := rows.Err(); err != nil {
		return []string{}, err
	}
	return devices, nil
}

func (AnR *AnalyticsRepository) GetCountriesForUrlById(ctx context.Context, UrlId int) ([]string, error) {
	var Countries []string
	GetDevicesQuery := `SELECT country FROM analytics WHERE url_id =$1;`
	rows, err := AnR.DB.DB.QueryContext(ctx, GetDevicesQuery, UrlId)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var country string
		err := rows.Scan(&country)
		if err != nil {
			return []string{}, err
		}
		Countries = append(Countries, country)
	}
	if err := rows.Err(); err != nil {
		return []string{}, err
	}
	return Countries, nil
}
