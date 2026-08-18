package handlers

import (
	"net/http"
	"url_shorter_gin/dto"
	"url_shorter_gin/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AS service.AuthService
}

func (AH *AuthHandler) RegisterHandler(c *gin.Context) {
	var registerData dto.RegisterRequest
	err := c.ShouldBindJSON(&registerData)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	ctx := c.Request.Context()
	userData, err := AH.AS.US.RegisterUser(ctx, registerData)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	ToDisplay := dto.DisplayUser{
		Id:       userData.Id,
		Role:     userData.Role,
		FullName: userData.FullName,
	}
	c.JSON(http.StatusCreated, ToDisplay)
}

func (AH *AuthHandler) LoginHandler(c *gin.Context) {
	var loginData dto.LoginRequest
	err := c.ShouldBindJSON(&loginData)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	ctx := c.Request.Context()
	loggedUser, err := AH.AS.US.Login(ctx, loginData)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	token, err := AH.AS.GenerateToken(loggedUser)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.Header("Authorization", token.Key)
}
