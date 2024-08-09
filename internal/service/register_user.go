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

type postRegisterSuccessResponse struct {
	UserID string `json:"user_id"`
}

func (s *Service) PostRegister(w http.ResponseWriter, r *http.Request) {
	var err error

	req, err := postRegisterParseBody(r.Body)
	if err != nil {
		s.logger.Errorf(err, "validation failed")

		err = WriteBadRequest(w)
		if err != nil {
			s.logger.Fatal("failed to send bad request response: %v", err)
		}

		return
	}

	err = postRegisterValidateRequest(req)
	if err != nil {
		s.logger.Errorf(err, "validation failed")

		err = WriteBadRequest(w)
		if err != nil {
			s.logger.Fatal("failed to send bad request response: %v", err)
		}

		return
	}

	resp, err := s.usecases.RegisterUser(r.Context(), usecases.RegisterUserRequest{
		Email:    string(*req.Email),
		Password: *req.Password,
		UserType: string(*req.UserType),
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

	err = postRegisterWriteSuccess(w, resp)
	if err != nil {
		s.logger.Fatal("failed to send successful response: %v", err)
	}
}

func postRegisterWriteSuccess(w http.ResponseWriter, useResp usecases.RegisterUserResponse) error {

	httpResp := postRegisterSuccessResponse{
		UserID: useResp.UUID.String(),
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

func postRegisterValidateRequest(req *api.PostRegisterJSONRequestBody) error {
	if req == nil {
		return errors.New("body must not be empty")
	}

	if *req.Email == "" {
		return errors.New("email must not be empty")
	}
	if *req.Password == "" {
		return errors.New("password must not be empty")
	}
	if *req.UserType == "" {
		return errors.New("user type must not be empty")
	}

	return nil
}

func postRegisterParseBody(body io.ReadCloser) (*api.PostRegisterJSONRequestBody, error) {
	var err error
	var apiRequest api.PostRegisterJSONRequestBody

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
