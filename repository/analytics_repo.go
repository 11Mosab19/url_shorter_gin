package repository

import (
	"context"
	"url_shorter_gin/database"
	"url_shorter_gin/models"
)

type AnalyticsRepository struct {
	DB *database.Database
}

func (AR *AnalyticsRepository) GetUrlAnalyticsByUrlId(ctx context.Context, UrlId int) ([]models.Analytics, error) {
	var analytics []models.Analytics
	GetAnalyticsQuery := `SELECT id,url_id,clicked_at,device_type,country FROM analytics WHERE url_id = $1 ORDER BY clicked_at DESC; `
	rows, err := AR.DB.DB.QueryContext(ctx, GetAnalyticsQuery, UrlId)
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
