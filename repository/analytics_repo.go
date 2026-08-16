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
	CreateAnalyticQuery := `INSERT INTO analytics (url_id,clicked_at,device_type,ip_address) VALUES ($1,$2,$3,$4) RETURNING id;`
	err := AnR.DB.DB.QueryRowContext(ctx, CreateAnalyticQuery, analytic.UrlID, analytic.ClickedAt, analytic.DeviceType, analytic.IpAddress).Scan(&analytic.Id)
	if err != nil {
		return models.Analytics{}, err
	}
	return analytic, nil
}

func (AnR *AnalyticsRepository) GetUrlAnalyticsByUrlId(ctx context.Context, UrlId int) ([]models.Analytics, error) {
	var analytics []models.Analytics
	GetAnalyticsQuery := `SELECT id,url_id,clicked_at,device_type,ip_address FROM analytics WHERE url_id = $1 ORDER BY clicked_at DESC; `
	rows, err := AnR.DB.DB.QueryContext(ctx, GetAnalyticsQuery, UrlId)
	if err != nil {
		return []models.Analytics{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var analytic models.Analytics
		err := rows.Scan(&analytic.Id, &analytic.UrlID, &analytic.ClickedAt, &analytic.DeviceType, &analytic.IpAddress)
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

func (AnR *AnalyticsRepository) GetIpAddressesForUrlById(ctx context.Context, UrlId int) ([]string, error) {
	var ipAddresses []string
	GetIpAddressQuery := `SELECT ip_address FROM analytics WHERE url_id =$1;`
	rows, err := AnR.DB.DB.QueryContext(ctx, GetIpAddressQuery, UrlId)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var ipAddress string
		err := rows.Scan(&ipAddress)
		if err != nil {
			return []string{}, err
		}
		ipAddresses = append(ipAddresses, ipAddress)
	}
	if err := rows.Err(); err != nil {
		return []string{}, err
	}
	return ipAddresses, nil
}
