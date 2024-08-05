package usecases

import (
	"fmt"
)

func (u *usecase) FlatUpdate(req FlatUpdateRequest) (FlatUpdateResponse, error) {

	flat, err := u.repo.FlatUpdate(req.FlatID, req.Status.String())
	if err != nil {
		return FlatUpdateResponse{}, fmt.Errorf("failed to update flat in repository: [%w]", err)
	}

	return FlatUpdateResponse{
		FlatID:   flat.FlatID,
		HouseID:  flat.HouseID,
		Price:    flat.Price,
		RoomsNum: flat.RoomsNum,
		Status:   flat.Status.String(),
	}, nil

}
