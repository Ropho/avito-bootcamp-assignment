package usecases

import (
	"context"
	"fmt"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/house"
)

func (u *usecase) HouseCreate(ctx context.Context, req HouseCreateRequest) (HouseCreateResponse, error) {
	var err error

	house := house.New(house.NewParams{
		Address:   req.Address,
		Year:      req.Year,
		Developer: req.Developer,
		Time:      u.time,
	})

	houseID, err := u.repo.HouseCreate(ctx, house)
	if err != nil {
		return HouseCreateResponse{}, fmt.Errorf("failed to create house in repository: [%w]", err)
	}

	house.HouseID = houseID

	return HouseCreateResponse{
		HouseID:   house.HouseID,
		Address:   house.Address,
		Year:      house.Year,
		Developer: house.Developer,
		CreatedAt: house.CreatedAt.String(),
		UpdatedAt: house.CreatedAt.String(),
	}, nil

}
