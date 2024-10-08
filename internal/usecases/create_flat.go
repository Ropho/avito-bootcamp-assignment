package usecases

import (
	"context"
	"fmt"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/flat"
)

func (u *usecases) FlatCreate(ctx context.Context, req FlatCreateRequest) (FlatCreateResponse, error) {
	var err error

	flat := flat.New(flat.NewParams{
		HouseID:  req.HouseID,
		Price:    req.Price,
		RoomsNum: req.RoomsNum,
	})

	flatID, err := u.repo.FlatCreate(ctx, flat)
	if err != nil {
		return FlatCreateResponse{}, fmt.Errorf("failed to create flat in repository: [%w]", err)
	}
	flat.FlatID = flatID

	err = u.SendEmailSubscribers(ctx, SendEmailSubscribersRequest{HouseID: int(req.HouseID)})
	if err != nil {
		return FlatCreateResponse{}, fmt.Errorf("failed to send notifications: [%w]", err)
	}

	return FlatCreateResponse{
		FlatID:   flat.FlatID,
		HouseID:  flat.HouseID,
		Price:    flat.Price,
		RoomsNum: flat.RoomsNum,
		Status:   flat.Status.String(),
	}, nil

}
