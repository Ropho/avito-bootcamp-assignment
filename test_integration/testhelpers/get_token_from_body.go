package testhelpers

import (
	"encoding/json"
	"fmt"
	"io"
)

func GetTokenFromBody(body io.ReadCloser) (string, error) {
	type Response struct {
		Token string `json:"token"`
	}

	var resp Response
	bytes, err := io.ReadAll(body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: [%w]", err)
	}
	defer body.Close()

	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal json struct: [%w]", err)
	}

	return resp.Token, nil
}
