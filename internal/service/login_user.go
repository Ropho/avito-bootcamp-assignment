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

type postLoginSuccessResponse struct {
	Token string `json:"token"`
}

func (s *Service) PostLogin(w http.ResponseWriter, r *http.Request) {
	var err error

	req, err := postLoginParseBody(r.Body)
	if err != nil {
		s.logger.Errorf(err, "validation failed")

		err = WriteBadRequest(w)
		if err != nil {
			s.logger.Fatal("failed to send bad request response: %v", err)
		}

		return
	}

	err = postLoginValidateRequest(req)
	if err != nil {
		s.logger.Errorf(err, "validation failed")

		err = WriteBadRequest(w)
		if err != nil {
			s.logger.Fatal("failed to send bad request response: %v", err)
		}

		return
	}

	resp, err := s.usecases.LoginUser(r.Context(), usecases.LoginUserRequest{
		UUID:     *req.Id,
		Password: *req.Password,
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

	err = postLoginWriteSuccess(w, resp)
	if err != nil {
		s.logger.Fatal("failed to send successful response: %v", err)
	}
}

func postLoginWriteSuccess(w http.ResponseWriter, useResp usecases.LoginUserResponse) error {

	httpResp := postLoginSuccessResponse{
		Token: useResp.Token,
	}

	bytes, err := json.Marshal(&httpResp)
	if err != nil {
		return fmt.Errorf("failed to marshal json object %v: [%v]", httpResp, err)
	}

	w.Header().Add(AuthorizationHeader, httpResp.Token)

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

func postLoginValidateRequest(req *api.PostLoginJSONRequestBody) error {
	if req == nil {
		return errors.New("body must not be empty")
	}

	if req.Id.String() == "" {
		return errors.New("user id must not be empty")
	}
	if *req.Password == "" {
		return errors.New("password must not be empty")
	}

	return nil
}

func postLoginParseBody(body io.ReadCloser) (*api.PostLoginJSONRequestBody, error) {
	var err error
	var apiRequest api.PostLoginJSONRequestBody

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
