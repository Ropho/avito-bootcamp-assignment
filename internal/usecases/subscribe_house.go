package usecases

import (
	"context"
	"fmt"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/subscription"
)

func (u *usecases) HouseSubscribe(ctx context.Context, req HouseSubscribeRequest) (HouseSubscribeResponse, error) {
	var err error

	sub := subscription.New(subscription.NewParams{
		UserID:  req.UserID,
		HouseID: req.HouseID,
		Email:   req.Email,
	})

	err = u.repo.SubscriptionCreate(ctx, sub)
	if err != nil {
		return HouseSubscribeResponse{}, fmt.Errorf("failed to create subscription in repository: [%w]", err)
	}

	return HouseSubscribeResponse{}, nil

}
