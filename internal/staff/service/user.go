package service

import (
	"context"
	"fmt"
	"projectsphere/eniqlo-store/internal/staff/entity"
	"projectsphere/eniqlo-store/internal/staff/repository"
	"projectsphere/eniqlo-store/pkg/middleware/auth"
	"projectsphere/eniqlo-store/pkg/protocol/msg"
	"projectsphere/eniqlo-store/pkg/validator"
)

type UserService struct {
	userRepo repository.UserRepo
	saltLen  int
	jwtAuth  auth.JWTAuth
}

func NewUserService(userRepo repository.UserRepo, saltLen int, jwtAuth auth.JWTAuth) UserService {
	return UserService{
		userRepo: userRepo,
		saltLen:  saltLen,
		jwtAuth:  jwtAuth,
	}
}

func (u UserService) Register(ctx context.Context, userParam *entity.UserParam) (entity.UserResponse, error) {
	if !validator.IsValidFullName(userParam.Name) {
		return entity.UserResponse{}, msg.BadRequest(msg.ErrInvalidFullName)
	}

	if !validator.IsEmailValid(userParam.Email) {
		return entity.UserResponse{}, msg.BadRequest(msg.ErrInvalidEmail)
	}

	if !validator.IsSolidPassword(userParam.Password) {
		return entity.UserResponse{}, msg.BadRequest(msg.ErrInvalidPassword)
	}

	if !validator.IsValidPhoneNumber(userParam.PhoneNumber) {
		return entity.UserResponse{}, msg.BadRequest(msg.ErrInvalidPhoneNumber)
	}

	userParam.Salt = auth.GenerateRandomAlphaNumeric(int(u.saltLen))
	hashedPassword := auth.GenerateHash([]byte(userParam.Password), []byte(userParam.Salt))
	userParam.Password = hashedPassword

	user, err := u.userRepo.CreateUser(ctx, *userParam)
	if err != nil {
		return entity.UserResponse{}, err
	}

	accessToken, err := u.jwtAuth.GenerateToken(user.UserId)
	if err != nil {
		return entity.UserResponse{}, err
	}

	return entity.UserResponse{
		UserId:      fmt.Sprint(user.UserId),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		AccessToken: accessToken,
	}, nil
}