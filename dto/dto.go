package dto

import "time"

type Device struct {
	Type        string
	TotalClicks int
}

type Country struct {
	Name        string
	TotalClicks int
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
	Id       int           `json:"id"`
	Role     string        `json:"role"`
	FullName string        `json:"full_name"`
	Urls     []ResponseUrl `json:"urls"`
}

type CreateURLRequest struct {
	OriginalUrl string     `json:"original_url" binding:"required,url"`
	ShortCode   string     `json:"short_code"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

type ResponseUrl struct {
	CreatorName string    `json:"creator_name"`
	OriginalUrl string    `json:"original_url"`
	ShortedCode string    `json:"shorted_code"`
	TotalClicks int       `json:"total_clicks"`
	DevicesUsed []Device  `json:"devices_used"`
	Countries   []Country `json:"countries"`
}

type UpdateURLRequest struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
	Expiration  string `json:"expiration"`
	Status      string `json:"status"`
}

type UpdateUserRequest struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
	FullName    string `json:"full_name"`
	Email       string `json:"email" binding:"email"`
}
