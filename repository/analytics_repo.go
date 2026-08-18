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
	CreateAnalyticQuery := `INSERT INTO analytics (url_id,clicked_at,device_type) VALUES ($1,$2,$3,$4) RETURNING id;`
	err := AnR.DB.DB.QueryRowContext(ctx, CreateAnalyticQuery, analytic.UrlID, analytic.ClickedAt, analytic.DeviceType).Scan(&analytic.Id)
	if err != nil {
		return models.Analytics{}, err
	}
	return analytic, nil
}

func (AnR *AnalyticsRepository) GetUrlAnalyticsByUrlId(ctx context.Context, UrlId int) ([]models.Analytics, error) {
	var analytics []models.Analytics
	GetAnalyticsQuery := `SELECT id,url_id,clicked_at,device_type FROM analytics WHERE url_id = $1 ORDER BY clicked_at DESC; `
	rows, err := AnR.DB.DB.QueryContext(ctx, GetAnalyticsQuery, UrlId)
	if err != nil {
		return []models.Analytics{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var analytic models.Analytics
		err := rows.Scan(&analytic.Id, &analytic.UrlID, &analytic.ClickedAt, &analytic.DeviceType)
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
	GetDevicesQuery := `SELECT device_type FROM analytics WHERE url_id =$1 ORDER BY clicked_at DESC;`
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

func (AnR *AnalyticsRepository) GetTotalClicksForUrlFromDevices(ctx context.Context, urlId int) ([]models.Device, error) {

	var devices []models.Device

	Query := `SELECT device_type,total_clicks FROM analytics JOIN urls ON
	analytics.url_id = urls.id 
	WHERE url_id = $1
	GROUP BY analytics.device_type;`

	rows, err := AnR.DB.DB.QueryContext(ctx, Query, urlId)

	if err != nil {
		return []models.Device{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var device models.Device
		err := rows.Scan(&device.Type, &device.TotalClicks)
		if err != nil {
			return []models.Device{}, err
		}
		devices = append(devices, device)
	}
	if err := rows.Err(); err != nil {
		return []models.Device{}, err
	}
	return devices, nil
}
