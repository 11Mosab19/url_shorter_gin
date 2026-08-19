package dto

import (
	"encoding/json"
	"time"
	"url_shorter_gin/models"
)

type Device struct {
	Type        string `json:"type"`
	TotalClicks int    `json:"total_clicks"`
}
type RegisterRequest struct {
	Email                string `json:"email" binding:"required,email"`
	Password             string `json:"password" binding:"required,min=8"`
	ConfirmationPassword string `json:"confirmation_password" binding:"required,eqfield=Password"`
	FullName             string `json:"full_name" binding:"required"`
}

type LoginRequest struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type DisplayUser struct {
	Id       int                  `json:"id"`
	Role     string               `json:"role"`
	FullName string               `json:"full_name"`
	Urls     []ResponseProfileUrl `json:"urls"`
}

type CreateURLRequest struct {
	OriginalUrl string     `json:"original_url" binding:"required,url"`
	ShortCode   string     `json:"short_code"`
	ExpiresAt   *time.Time `json:"expires_at"`
	Password    string     `json:"password" binding:"min=8"`
}

type ResponseUrl struct {
	Id          int             `json:"id"`
	OriginalUrl string          `json:"original_url"`
	ShortedCode string          `json:"shorted_code"`
	TotalClicks int             `json:"total_clicks"`
	Status      string          `json:"status"`
	DevicesUsed []models.Device `json:"devices_used"`
}

type UpdateURLRequest struct {
	NewPassword string          `json:"new_password" binding:"min=8"`
	OldPassword string          `json:"old_password"`
	Expiration  json.RawMessage `json:"expiration"`
	Status      string          `json:"status"`
}

type UpdateUserRequest struct {
	NewPassword string `json:"new_password" binding:"min=8"`
	OldPassword string `json:"old_password"`
	FullName    string `json:"full_name"`
	Email       string `json:"email" binding:"email"`
}

type SetPasswordRequest struct {
	NewPassword          string `json:"new_password" binding:"min=8"`
	ConfirmationPassword string `json:"confirmation_password" binding:"required,eqfield=NewPassword"`
}

type Token struct {
	Key string `json:"token"`
}
type ResponseProfileUrl struct {
	OriginalUrl string `json:"original_url"`
	ShortedCode string `json:"shorted_code"`
	TotalClicks int    `json:"total_clicks"`
	Status      string `json:"status"`
}

type UrlPasswordReq struct {
	Password string `form:"password"`
}
