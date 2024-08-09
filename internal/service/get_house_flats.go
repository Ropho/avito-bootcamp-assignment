package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Ropho/avito-bootcamp-assignment/api"
	"github.com/Ropho/avito-bootcamp-assignment/internal/usecases"
)

type getHouseIdSuccessResponse struct {
	Flats []getHouseIdFlat `json:"flats"`
}
type getHouseIdFlat struct {
	FlatID   uint32 `json:"id"`
	HouseID  uint32 `json:"house_id"`
	Price    uint32 `json:"price"`
	RoomsNum uint32 `json:"rooms"`
	Status   string `json:"status"`
}

func (s *Service) GetHouseId(w http.ResponseWriter, r *http.Request, id api.HouseId) {
	var err error

	err = getHouseIdValidateRequest(id)
	if err != nil {
		s.logger.Errorf(err, "validation failed")

		err = WriteBadRequest(w)
		if err != nil {
			s.logger.Fatal("failed to send bad request response: %v", err)
		}

		return
	}

	resp, err := s.usecases.GetHouseFlats(r.Context(), usecases.GetHouseFlatsRequest{
		HouseID:      uint32(id),
		OnlyApproved: isOnlyApprovedFlats(r.Context()),
	})
	if err != nil {
		s.logger.Errorf(err, "service logic failed")

		requestID, ok := r.Context().Value(RequestIDKey{}).(string)
		if !ok {
			s.logger.Fatal("failed to conver request id to string: %v", err)
		}
		err = WriteInternal(w, requestID)
		if err != nil {
			s.logger.Fatal("failed to send internal response: %v", err)
		}

		return
	}
	fmt.Println(resp)
	err = getHouseIdWriteSuccess(w, resp)
	if err != nil {
		s.logger.Fatal("failed to send successful response: %v", err)
	}
}

func isOnlyApprovedFlats(ctx context.Context) bool {
	userType := ctx.Value(UserTypeKey{})

	fmt.Println(userType)

	userTypeString, ok := userType.(string)
	if !ok {
		return false
	}

	return userTypeString != string(api.Moderator)
}

func getHouseIdWriteSuccess(w http.ResponseWriter, useResp usecases.GetHouseFlatsResponse) error {

	httpResp := getHouseIdSuccessResponse{
		Flats: make([]getHouseIdFlat, len(useResp.Flats)),
	}
	for index := range useResp.Flats {
		httpResp.Flats[index].FlatID = useResp.Flats[index].FlatID
		httpResp.Flats[index].HouseID = useResp.Flats[index].HouseID
		httpResp.Flats[index].Price = useResp.Flats[index].Price
		httpResp.Flats[index].RoomsNum = useResp.Flats[index].Rooms
		httpResp.Flats[index].Status = useResp.Flats[index].Status
	}

	bytes, err := json.Marshal(&httpResp)
	if err != nil {
		return fmt.Errorf("failed to marshal json object %v: [%v]", httpResp, err)
	}
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

func getHouseIdValidateRequest(id api.HouseId) error {
	if id < 1 {
		return errors.New("house id must be greater than 0")
	}

	return nil
}
