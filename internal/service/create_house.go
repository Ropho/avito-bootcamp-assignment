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

type postHouseCreateSuccessResponse struct {
	HouseID   uint32 `json:"id"`
	Address   string `json:"address"`
	Year      uint32 `json:"year"`
	Developer string `json:"developer"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"update_at"`
}

func (s *Service) PostHouseCreate(w http.ResponseWriter, r *http.Request) {
	var err error

	req, err := postHouseCreateParseBody(r.Body)
	if err != nil {
		s.logger.Errorf(err, "validation failed")

		err = WriteBadRequest(w)
		if err != nil {
			s.logger.Fatal("failed to send bad request response: %v", err)
		}

		return
	}

	err = postHouseCreateValidateRequest(req)
	if err != nil {
		s.logger.Errorf(err, "validation failed")

		err = WriteBadRequest(w)
		if err != nil {
			s.logger.Fatal("failed to send bad request response: %v", err)
		}

		return
	}

	developer := ""
	if req.Developer != nil {
		developer = *req.Developer
	}
	resp, err := s.usecases.HouseCreate(r.Context(), usecases.HouseCreateRequest{
		Address:   req.Address,
		Year:      uint32(req.Year),
		Developer: developer,
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

	err = postHouseCreateWriteSuccess(w, resp)
	if err != nil {
		s.logger.Fatal("failed to send successful response: %v", err)
	}

}

func postHouseCreateWriteSuccess(w http.ResponseWriter, useResp usecases.HouseCreateResponse) error {

	httpResp := postHouseCreateSuccessResponse{
		HouseID:   useResp.HouseID,
		Address:   useResp.Address,
		Year:      useResp.Year,
		Developer: useResp.Developer,
		CreatedAt: useResp.CreatedAt,
		UpdatedAt: useResp.UpdatedAt,
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

func postHouseCreateValidateRequest(req *api.PostHouseCreateJSONRequestBody) error {
	if req == nil {
		return errors.New("body must not be empty")
	}

	if req.Address == " " {
		return errors.New("house address must not be empty")
	}
	if req.Year < 1 {
		return errors.New("year must be greater than 0")
	}
	// req.Developer may be empty

	return nil
}

func postHouseCreateParseBody(body io.ReadCloser) (*api.PostHouseCreateJSONRequestBody, error) {
	var err error
	var apiRequest api.PostHouseCreateJSONRequestBody

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
