package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Id             int
	Email          string
	HashedPassword string
	FullName       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Role           string
}

type Url struct {
	Id             int
	Status         string
	OriginalUrl    string
	ShortCode      string
	CreatedAt      time.Time
	ExpiresAt      *time.Time
	HashedPassword string
	UpdatedAt      time.Time
	UserID         int
	TotalClicks    int
}

type Analytics struct {
	Id         int
	UrlID      int
	ClickedAt  time.Time
	DeviceType string
	UserAgent  string
}

type Claims struct {
	UserId int
	Role   string
	jwt.RegisteredClaims
}
