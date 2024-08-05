package usecases

import (
	"fmt"
)

func (u *usecase) GetHouseFlats(req GetHouseFlatsRequest) (GetHouseFlatsResponse, error) {
	var err error

	flats, err := u.repo.GetHouseFlats(req.HouseID, req.OnlyApproved)
	if err != nil {
		return GetHouseFlatsResponse{}, fmt.Errorf("failed to get flats in the house with id %d in repository: [%w]", req.HouseID, err)
	}

	resp := GetHouseFlatsResponse{
		Flats: make([]struct {
			FlatID  uint32
			HouseID uint32
			Price   uint32
			Rooms   uint32
			Status  string
		}, len(flats)),
	}

	for index := range flats {
		resp.Flats[index].FlatID = flats[index].FlatID
		resp.Flats[index].HouseID = flats[index].HouseID
		resp.Flats[index].Price = flats[index].Price
		resp.Flats[index].Rooms = flats[index].RoomsNum
		resp.Flats[index].Status = flats[index].Status.String()
	}

	return resp, nil
}
