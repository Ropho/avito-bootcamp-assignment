package usecases

import (
	"context"
	"fmt"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/hash"
)

func (u *usecases) LoginUser(ctx context.Context, req LoginUserRequest) (LoginUserResponse, error) {
	var err error

	user, err := u.repo.GetUserByID(ctx, req.UUID)
	if err != nil {
		return LoginUserResponse{}, fmt.Errorf("failed to get user by id in repository: [%w]", err)
	}

	if !hash.IsValidPassWithSalt(user.EncryptedPassword, req.Password, user.Salt) {
		return LoginUserResponse{}, fmt.Errorf("incorrect password")
	}

	token, err := u.jwtService.GenerateUserAccessJWT(user.ID.String(), user.Type.String())
	if err != nil {
		return LoginUserResponse{}, fmt.Errorf("failed to generate token: [%w]", err)
	}

	return LoginUserResponse{
		Token: token,
	}, nil
}
