package grpcserver

import (
	"context"

	"github.com/Garvit-Jethwani/user-management-service/database"
	"github.com/Garvit-Jethwani/user-management-service/models"
)

type UserServiceServer struct {
	UnimplementedUserServiceServer
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *GetUserRequest) (*UserResponse, error) {
	user, err := database.GetUserByID(req.UserId)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := database.CreateUser(user); err != nil {
		return nil, err
	}

	return &UserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
