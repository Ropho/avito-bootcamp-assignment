package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (u *usecases) GetDummyLogin(ctx context.Context, req GetDummyLoginRequest) (GetDummyLoginResponse, error) {

	userID, err := uuid.NewV6()
	if err != nil {
		return GetDummyLoginResponse{}, fmt.Errorf("failed to generate dummy user id: %w", err)
	}

	token, err := u.jwtService.GenerateUserAccessJWT(userID.String(), req.UserType)
	if err != nil {
		return GetDummyLoginResponse{}, fmt.Errorf("failed to generate token: [%w]", err)
	}

	return GetDummyLoginResponse{
		Token: token,
	}, nil
}
