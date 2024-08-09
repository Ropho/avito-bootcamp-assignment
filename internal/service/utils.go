package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	InvlaidRequestMessage    = "Невалидные данные ввода"
	InternalErrorMessage     = "Ошибка сервера"
	UnauthorizedErrorMessage = "Неавторизованный доступ"
)

const (
	FlatUpdateRequestMessage = "Успещно обновлена квартира"
)

const RetryHeader = "Retry-After"
const RetryHeaderValue = 30

type JSON500 struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Code      int    `json:"code"`
}

func WriteInternal(w http.ResponseWriter, requestID string) error {

	resp := JSON500{
		Message:   InternalErrorMessage,
		RequestID: requestID,
		Code:      http.StatusInternalServerError,
	}

	bytes, err := json.Marshal(&resp)
	if err != nil {
		return fmt.Errorf("failed to marshal json object %v: [%v]", resp, err)
	}
	w.WriteHeader(http.StatusInternalServerError)

	w.Header().Add(RetryHeader, strconv.FormatInt(RetryHeaderValue, 10))

	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

func WriteBadRequest(w http.ResponseWriter) error {

	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte(InvlaidRequestMessage))
	if err != nil {
		return fmt.Errorf("failed to write response: [%w]", err)
	}

	return nil
}

func WriteUnauthorized(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusUnauthorized)
	_, err := w.Write([]byte(UnauthorizedErrorMessage))
	if err != nil {
		return fmt.Errorf("failed to write response: [%w]", err)
	}

	return nil
}
