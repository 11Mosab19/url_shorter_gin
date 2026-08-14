package service

import (
	"context"
	"url_shorter_gin/apperrors"
	"url_shorter_gin/dto"
	"url_shorter_gin/models"
	"url_shorter_gin/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func (US *UserService) RegisterUser(ctx context.Context, registerData dto.RegisterRequest) (models.User, error) {
	var user models.User
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, apperrors.HashingErr
	}
	registerData.Password = string(HashedPassword)

	user.Email = registerData.Email
	user.FullName = registerData.FullName
	user.HashedPassword = registerData.Password

	user, err = US.repo.CreateUser(ctx, user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (US *UserService) Login(ctx context.Context, loginData dto.LoginRequest) (models.User, error) {
	user, err := US.repo.GetUserByEmail(ctx, loginData.Email)
	if err != nil {
		return models.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(loginData.Password))
	if err != nil {
		return models.User{}, apperrors.ErrWrongUserPassword
	}
	return user, nil
}

func (US *UserService) GetUserById(ctx context.Context, id int) (models.User, error) {
	return US.repo.GetUserById(ctx, id)
}

func (US *UserService) UpdateUserPassword(ctx context.Context, req dto.UpdateUserRequest, id int) (models.User, error) {
	user, err := US.GetUserById(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.OldPassword))
	if err != nil {
		return models.User{}, apperrors.ErrWrongUserPassword
	}
	Hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, apperrors.HashingErr
	}
	user, err = US.repo.UpdateUserHashedPassword(ctx, string(Hashed), id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (US *UserService) UpdateUserEmail(ctx context.Context, req dto.UpdateUserRequest, id int) (models.User, error) {
	user, err := US.repo.UpdateUserEmail(ctx, id, req.Email)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (US *UserService) UpdateUserFullName(ctx context.Context, req dto.UpdateUserRequest, id int) (models.User, error) {
	user, err := US.repo.UpdateUserFullName(ctx, id, req.FullName)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
