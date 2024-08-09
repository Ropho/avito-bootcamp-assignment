package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"

	"github.com/Ropho/avito-bootcamp-assignment/api"
	"github.com/Ropho/avito-bootcamp-assignment/internal/usecases"
)

var successMessage = "Успешно оформлена подписка"

func (s *Service) PostHouseIdSubscribe(w http.ResponseWriter, r *http.Request, id api.HouseId) {
	var err error

	userID, ok := r.Context().Value(UserIDKey{}).(string)
	if !ok {
		s.logger.Fatal("failed to get userID from context: %v", err)
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.Errorf(err, "validation failed")

		err = WriteBadRequest(w)
		if err != nil {
			s.logger.Fatal("failed to send bad request response: %v", err)
		}

		return
	}

	req, err := postHouseIdSubscribeParseBody(r.Body)
	if err != nil {
		s.logger.Errorf(err, "validation failed")

		err = WriteBadRequest(w)
		if err != nil {
			s.logger.Fatal("failed to send bad request response: %v", err)
		}

		return
	}

	err = postHouseIdSubscribeValidateRequest(req)
	if err != nil {
		s.logger.Errorf(err, "validation failed")

		err = WriteBadRequest(w)
		if err != nil {
			s.logger.Fatal("failed to send bad request response: %v", err)
		}

		return
	}

	_, err = s.usecases.HouseSubscribe(r.Context(), usecases.HouseSubscribeRequest{
		Email:   string(req.Email),
		UserID:  userUUID,
		HouseID: id,
	})
	if err != nil {
		s.logger.Errorf(err, "service logic failed")

		requestID, ok := r.Context().Value(RequestIDKey{}).(string)
		if !ok {
			s.logger.Fatal("failed to convert request id to string: %v", err)
		}
		err = WriteInternal(w, requestID)
		if err != nil {
			s.logger.Fatal("failed to send internal response: %v", err)
		}
		return
	}

	err = postHouseIdSubscribeWriteSuccess(w)
	if err != nil {
		s.logger.Fatal("failed to send successful response: %v", err)
	}
}

func postHouseIdSubscribeWriteSuccess(w http.ResponseWriter) error {

	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(successMessage))
	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

func postHouseIdSubscribeValidateRequest(req *api.PostHouseIdSubscribeJSONRequestBody) error {
	if req == nil {
		return errors.New("body must not be empty")
	}
	if req.Email == "" {
		return errors.New("email must not be empty")
	}

	return nil
}

func postHouseIdSubscribeParseBody(body io.ReadCloser) (*api.PostHouseIdSubscribeJSONRequestBody, error) {
	var err error
	var apiRequest api.PostHouseIdSubscribeJSONRequestBody

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
