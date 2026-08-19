package handlers

import (
	"net/http"
	"url_shorter_gin/apperrors"
	"url_shorter_gin/dto"
	"url_shorter_gin/models"
	"url_shorter_gin/service"

	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	US  service.UserService
	UrS service.UrlService
}

func (UH *UsersHandler) GetUserProfile(c *gin.Context) {
	ctx := c.Request.Context()

	UserId, exist := c.Get("userId")
	if !exist {
		c.Error(apperrors.ErrUnauthorized)
		c.Abort()
		return
	}

	userData, err := UH.US.GetUserById(ctx, UserId.(int))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	urls, err := UH.UrS.GetUserUrls(ctx, userData.Id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	var ResponseProfileUrls []dto.ResponseProfileUrl

	for _, url := range urls {
		dtoUrl := dto.ResponseProfileUrl{
			OriginalUrl: url.OriginalUrl,
			ShortedCode: url.ShortCode,
			TotalClicks: url.TotalClicks,
			Status:      url.Status,
		}
		ResponseProfileUrls = append(ResponseProfileUrls, dtoUrl)
	}

	DisplayUser := dto.DisplayUser{
		Id:       userData.Id,
		Role:     userData.Role,
		Urls:     ResponseProfileUrls,
		FullName: userData.FullName,
	}

	c.JSON(http.StatusOK, DisplayUser)
}

func (UH *UsersHandler) UpdateUser(c *gin.Context) {
	var req dto.UpdateUserRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	UserId, exist := c.Get("userId")
	if !exist {
		c.Error(apperrors.ErrUnauthorized)
		c.Abort()
		return
	}

	userID, ok := UserId.(int)
	if !ok {
		c.Error(apperrors.ErrUnauthorized)
		c.Abort()
		return
	}

	ctx := c.Request.Context()

	var user models.User

	if req.Email != "" {
		user, err = UH.US.UpdateUserEmail(ctx, req, userID)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
	}

	if req.FullName != "" {
		user, err = UH.US.UpdateUserFullName(ctx, req, userID)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
	}

	if req.NewPassword != "" {
		user, err = UH.US.UpdateUserPassword(ctx, req, userID)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
	}

	user, err = UH.US.GetUserById(ctx, userID)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, dto.DisplayUser{
		Id:       user.Id,
		Role:     user.Role,
		FullName: user.FullName,
	})
}
