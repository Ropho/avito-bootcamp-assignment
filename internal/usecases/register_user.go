package usecases

import (
	"context"
	"fmt"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/user"
)

func (u *usecases) RegisterUser(ctx context.Context, req RegisterUserRequest) (RegisterUserResponse, error) {
	var err error

	user, err := user.New(user.NewParams{
		Email:    req.Email,
		Password: req.Password,
		UserType: req.UserType,
	})
	if err != nil {
		return RegisterUserResponse{}, fmt.Errorf("failed to create user model: [%w]", err)
	}

	err = u.repo.RegisterUser(ctx, user)
	if err != nil {
		return RegisterUserResponse{}, fmt.Errorf("failed to create user in repository: [%w]", err)
	}

	return RegisterUserResponse{
		UUID: user.ID,
	}, nil
}
