package handlers

import (
	"net/http"
	"strconv"
	"url_shorter_gin/apperrors"
	"url_shorter_gin/dto"
	"url_shorter_gin/service"

	"github.com/gin-gonic/gin"
)

type UrlsHandler struct {
	UrS service.UrlService
	AH  service.AnalyticsService
}

func (UrH *UrlsHandler) GetUrlById(c *gin.Context) {
	ctx := c.Request.Context()

	urlId := c.Param("id")

	userId, exists := c.Get("userId")

	if !exists {
		c.Error(apperrors.ErrUnauthorized)
		c.Abort()
		return
	}

	UserId, ok := userId.(int)

	if !ok {
		c.Error(apperrors.ErrUnauthorized)
		c.Abort()
		return
	}

	UrlId, err := strconv.Atoi(urlId)

	if err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	url, err := UrH.UrS.GetUrlById(ctx, UrlId, UserId)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	devices, err := UrH.AH.GetUrlUsersDevices(ctx, UrlId)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	ToDisplayUrl := dto.ResponseUrl{
		OriginalUrl: url.OriginalUrl,
		ShortedCode: url.ShortCode,
		TotalClicks: url.TotalClicks,
		DevicesUsed: devices,
	}

	c.JSON(http.StatusOK, ToDisplayUrl)
}

/*func (UrH *UrlsHandler) CreateUrl(c *gin.Context) {
	var createUrlReq dto.CreateURLRequest

	ctx := c.Request.Context()

	urlId := c.Param("id")

	err := c.ShouldBindJSON(&createUrlReq)

	if err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	UrlId, err := strconv.Atoi(urlId)

	if err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	url, err := UrH.UrS.CreateUrl(ctx)

}*/
