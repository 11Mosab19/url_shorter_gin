package apperrors

import "errors"

var (
	ErrUrlNotFound     = errors.New("url not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailIsUsed     = errors.New("this email is used")
	ErrShortCodeUsed   = errors.New("this short code is used")
	ErrAlreadyExistUrl = errors.New("this user has this url used")
)
