package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Ropho/avito-bootcamp-assignment/api"
	"github.com/Ropho/avito-bootcamp-assignment/internal/usecases"
)

type postFlatCreateSuccessResponse struct {
	FlatID  uint32 `json:"id"`
	HouseID uint32 `json:"house_id"`
	Price   uint32 `json:"price"`
	Rooms   uint32 `json:"rooms"`
	Status  string `json:"status"`
}

func (s *Service) PostFlatCreate(w http.ResponseWriter, r *http.Request) {
	var err error

	req, err := postFlatCreateParseBody(r.Body)
	if err != nil {
		s.logger.Errorf(err, "validation failed")

		err = WriteBadRequest(w)
		if err != nil {
			s.logger.Fatal("failed to send bad request response: %v", err)
		}

		return
	}

	err = postFlatCreateValidateRequest(req)
	if err != nil {
		s.logger.Errorf(err, "validation failed")

		err = WriteBadRequest(w)
		if err != nil {
			s.logger.Fatal("failed to send bad request response: %v", err)
		}

		return
	}

	resp, err := s.usecases.FlatCreate(r.Context(), usecases.FlatCreateRequest{
		HouseID:  uint32(req.HouseId),
		Price:    uint32(req.Price),
		RoomsNum: uint32(*req.Rooms),
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

	err = postFlatCreateWriteSuccess(w, resp)
	if err != nil {
		s.logger.Fatal("failed to send successful response: %v", err)
	}
}

func postFlatCreateWriteSuccess(w http.ResponseWriter, useResp usecases.FlatCreateResponse) error {

	httpResp := postFlatCreateSuccessResponse{
		FlatID:  useResp.FlatID,
		HouseID: useResp.HouseID,
		Price:   useResp.Price,
		Rooms:   useResp.RoomsNum,
		Status:  useResp.Status,
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

func postFlatCreateValidateRequest(req *api.PostFlatCreateJSONRequestBody) error {
	if req == nil {
		return errors.New("body must not be empty")
	}

	if req.HouseId < 1 {
		return errors.New("house id must be grater than 0")
	}
	if req.Price < 0 {
		return errors.New("flat price must be greater or equal than 0")
	}
	if req.Rooms == nil {
		return errors.New("rooms number must be present")
	}

	if *req.Rooms < 1 {
		return errors.New("rooms number must be greater than 0")
	}

	return nil
}

func postFlatCreateParseBody(body io.ReadCloser) (*api.PostFlatCreateJSONRequestBody, error) {
	var err error
	var apiRequest api.PostFlatCreateJSONRequestBody

	bytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: [%w]", err)
	}
	defer body.Close()

	err = json.Unmarshal(bytes, &apiRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json struct: [%w]", err)
	}

	return &apiRequest, nil
}
