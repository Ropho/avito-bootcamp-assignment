package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ropho/avito-bootcamp-assignment/api"
	"github.com/Ropho/avito-bootcamp-assignment/internal/usecases"
)

type getDummyLoginSuccessResponse struct {
	Token string `json:"token"`
}

func (s *Service) GetDummyLogin(w http.ResponseWriter, r *http.Request, params api.GetDummyLoginParams) {
	var err error

	err = getDummyLoginValidateRequest(params)
	if err != nil {
		s.logger.Errorf(err, "validation failed")

		err = WriteBadRequest(w)
		if err != nil {
			s.logger.Fatal("failed to send bad request response: %v", err)
		}

		return
	}

	resp, err := s.usecases.GetDummyLogin(r.Context(), usecases.GetDummyLoginRequest{
		UserType: string(params.UserType),
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

	err = getDummyLoginWriteSuccess(w, resp)
	if err != nil {
		s.logger.Fatal("failed to send successful response: %v", err)
	}
}

func getDummyLoginWriteSuccess(w http.ResponseWriter, useResp usecases.GetDummyLoginResponse) error {

	httpResp := getDummyLoginSuccessResponse{
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

func getDummyLoginValidateRequest(req api.GetDummyLoginParams) error {

	if !(req.UserType == api.Moderator || req.UserType == api.Client) {
		return fmt.Errorf("unknown user type %s", req.UserType)
	}

	return nil
}
