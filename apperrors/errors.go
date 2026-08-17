package apperrors

import (
	"net/http"
)

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrUrlNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "url not found",
	}

	ErrUserNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "user not found",
	}

	ErrEmailIsUsed = &AppError{
		Code:    http.StatusConflict,
		Message: "this email is used",
	}

	ErrShortCodeUsed = &AppError{
		Code:    http.StatusConflict,
		Message: "this short code is used",
	}

	ErrAlreadyExistUrl = &AppError{
		Code:    http.StatusConflict,
		Message: "this user has this url used",
	}

	ErrNoUrlWithId = &AppError{
		Code:    http.StatusNotFound,
		Message: "no url with this id",
	}

	ErrWrongUrlPassword = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "wrong url password",
	}

	ErrWrongUserPassword = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "wrong user password",
	}

	HashingErr = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "error while hashing",
	}

	ErrGenerateShortCode = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "couldn't generate shortcode",
	}

	ErrInvalidStatusInput = &AppError{
		Code:    http.StatusBadRequest,
		Message: "invalid update statue",
	}

	ErrBrokenUrl = &AppError{
		Code:    http.StatusForbidden,
		Message: "this url isn't active",
	}

	ErrExpiredUrl = &AppError{
		Code:    http.StatusGone,
		Message: "url expired",
	}

	ErrUnauthorized = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "not allowed",
	}

	ErrInvalidExpirationInput = &AppError{
		Code:    http.StatusBadRequest,
		Message: "invalid expiration date",
	}

	ErrWrongConfirmationPassword = &AppError{
		Code:    http.StatusBadRequest,
		Message: "wrong confirmation password",
	}

	ErrNotSupportedAuthenticateMethod = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "not supported",
	}
)
