package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Ropho/avito-bootcamp-assignment/api"
	"github.com/Ropho/avito-bootcamp-assignment/internal/models/flat"
	"github.com/Ropho/avito-bootcamp-assignment/internal/usecases"
)

type postFlatUpdateSuccessResponse struct {
	ID      uint32 `json:"id"`
	HouseID uint32 `json:"house_id"`
	Price   uint32 `json:"price"`
	Rooms   uint32 `json:"rooms"`
	Status  string `json:"status"`
}

func (s *Service) PostFlatUpdate(w http.ResponseWriter, r *http.Request) {
	var err error

	req, err := postFlatUpdateParseBody(r.Body)
	if err != nil {
		s.logger.Errorf(err, "validation failed")

		err = WriteBadRequest(w)
		if err != nil {
			s.logger.Fatal("failed to send bad request response: %v", err)
		}

		return
	}

	err = postFlatUpdateValidateRequest(req)
	if err != nil {
		s.logger.Errorf(err, "validation failed")

		err = WriteBadRequest(w)
		if err != nil {
			s.logger.Fatal("failed to send bad request response: %v", err)
		}

		return
	}

	resp, err := s.usecases.FlatUpdate(r.Context(), usecases.FlatUpdateRequest{
		FlatID: uint32(req.Id),
		Status: flat.GetStatusFromString(string(*req.Status)),
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

	err = postFlatUpdateWriteSuccess(w, resp)
	if err != nil {
		s.logger.Fatal("failed to send successful response: %v", err)
	}

}

func postFlatUpdateWriteSuccess(w http.ResponseWriter, useResp usecases.FlatUpdateResponse) error {

	httpResp := postFlatUpdateSuccessResponse{
		ID:      useResp.FlatID,
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

func postFlatUpdateValidateRequest(req *api.PostFlatUpdateJSONRequestBody) error {
	if req == nil {
		return errors.New("body must not be empty")
	}
	if req.Id < 1 {
		return errors.New("flat id must be greater than 0")
	}

	return nil
}

func postFlatUpdateParseBody(body io.ReadCloser) (*api.PostFlatUpdateJSONRequestBody, error) {
	var err error
	var apiRequest api.PostFlatUpdateJSONRequestBody

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
