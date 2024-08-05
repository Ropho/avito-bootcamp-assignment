package usecases

import (
	"fmt"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/flat"
)

func (u *usecase) FlatCreate(req FlatCreateRequest) (FlatCreateResponse, error) {
	var err error

	flat := flat.New(flat.NewParams{
		HouseID:  req.HouseID,
		Price:    req.Price,
		RoomsNum: req.RoomsNum,
	})

	flatID, err := u.repo.FlatCreate(flat)
	if err != nil {
		return FlatCreateResponse{}, fmt.Errorf("failed to create flat in repository: [%w]", err)
	}
	flat.FlatID = flatID

	return FlatCreateResponse{
		FlatID:   flat.FlatID,
		HouseID:  flat.HouseID,
		Price:    flat.Price,
		RoomsNum: flat.RoomsNum,
		Status:   flat.Status.String(),
	}, nil

}
