package handlers

import (
	"net/http"
	"strconv"
	"url_shorter_gin/apperrors"
	"url_shorter_gin/dto"
	"url_shorter_gin/models"
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

	ToDisplayUrl := dto.ResponseProfileUrl{
		OriginalUrl: url.OriginalUrl,
		ShortedCode: url.ShortCode,
		TotalClicks: url.TotalClicks,
	}

	c.JSON(http.StatusOK, ToDisplayUrl)
}

func (UrH *UrlsHandler) CreateUrl(c *gin.Context) {
	var createUrlReq dto.CreateURLRequest

	ctx := c.Request.Context()

	err := c.ShouldBindJSON(&createUrlReq)

	if err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

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

	url, err := UrH.UrS.CreateUrl(ctx, createUrlReq, UserId)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	resUrl := dto.ResponseUrl{
		Id:          url.Id,
		ShortedCode: url.ShortCode,
		OriginalUrl: url.OriginalUrl,
		Status:      url.Status,
		TotalClicks: url.TotalClicks,
	}

	c.JSON(http.StatusCreated, resUrl)

}

func (UrH *UrlsHandler) RedirectByShortCode(c *gin.Context) {

	ctx := c.Request.Context()

	shortCode := c.Param("shortCode")

	url, err := UrH.UrS.GetUrlByShortcode(ctx, shortCode)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	if url.HashedPassword != "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "password required"})
		return
	}

	url, err = UrH.UrS.RedirectByShortCode(ctx, shortCode, "")

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	analytic := models.Analytics{
		UrlID:     url.Id,
		UserAgent: c.GetHeader("User-Agent"),
	}

	_, err = UrH.AH.CreateAnalytic(ctx, analytic)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Redirect(http.StatusPermanentRedirect, url.OriginalUrl)
}

func (UrH *UrlsHandler) RedirectVerifiedUrlByShortcode(c *gin.Context) {
	ctx := c.Request.Context()

	shortCode := c.Param("shortCode")

	var password dto.UrlPasswordReq

	err := c.ShouldBind(&password)

	if err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	url, err := UrH.UrS.GetUrlByShortcode(ctx, shortCode)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	url, err = UrH.UrS.RedirectByShortCode(ctx, shortCode, password.Password)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	analytic := models.Analytics{
		UrlID:     url.Id,
		UserAgent: c.GetHeader("User-Agent"),
	}

	_, err = UrH.AH.CreateAnalytic(ctx, analytic)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Redirect(http.StatusPermanentRedirect, url.OriginalUrl)
}

func (UrH *UrlsHandler) Delete(c *gin.Context) {
	UserId, exists := c.Get("userId")

	if !exists {
		c.Error(apperrors.ErrUnauthorized)
		c.Abort()
		return
	}

	userId, ok := UserId.(int)

	if !ok {
		c.Error(apperrors.ErrUnauthorized)
		c.Abort()
		return
	}

	ctx := c.Request.Context()

	UrlId := c.Param("id")

	urlId, err := strconv.Atoi(UrlId)

	if err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	err = UrH.UrS.DeleteUrlById(ctx, urlId, userId)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}

func (UrH *UrlsHandler) UpdateUrl(c *gin.Context) {

	UserId, exists := c.Get("userId")

	if !exists {
		c.Error(apperrors.ErrUnauthorized)
		c.Abort()
		return
	}

	userId, ok := UserId.(int)

	if !ok {
		c.Error(apperrors.ErrUnauthorized)
		c.Abort()
		return
	}

	ctx := c.Request.Context()

	UrlId := c.Param("id")

	urlId, err := strconv.Atoi(UrlId)

	if err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	var updateReq dto.UpdateURLRequest

	err = c.ShouldBindJSON(&updateReq)

	if err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	var url models.Url

	if updateReq.Status != "" {
		url, err = UrH.UrS.UpdateUrlStatus(ctx, userId, updateReq, urlId)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
	}

	if updateReq.NewPassword != "" {
		url, err = UrH.UrS.UpdateUrlPassword(ctx, updateReq, urlId, userId)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
	}

	if updateReq.Expiration != nil {
		url, err = UrH.UrS.UpdateUrlExpiration(ctx, updateReq, urlId, userId)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
	}

	url, err = UrH.UrS.GetUrlById(ctx, urlId, userId)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, url)
}
