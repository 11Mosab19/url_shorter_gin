package handlers

import "url_shorter_gin/service"

type AuthHandler struct {
	AS service.AuthService
}
