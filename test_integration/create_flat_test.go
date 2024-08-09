package testintegration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/require"
)

type FlatCreateSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *FlatCreateSuite) BeforeAll(provider.T) {}

func (s *FlatCreateSuite) TestService_FlatCreate(t provider.T) {

	type bodyFlatCreateType struct {
		HouseID int32 `json:"house_id"`
		Price   int32 `json:"price"`
		Rooms   int32 `json:"rooms"`
	}

	resp, err := http.Get("http://localhost:8080/dummyLogin?user_type=moderator")
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)

	token, err := getTokenFromBody(resp.Body)
	require.NoError(t, err)

	body := bodyFlatCreateType{
		HouseID: 1,
		Price:   1000,
		Rooms:   12,
	}
	bodyBytes, err := json.Marshal(&body)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "http://localhost:8080/flat/create", bytes.NewReader(bodyBytes))
	require.NoError(t, err)
	req.Header.Set("Authorization", token)

	client := http.Client{}
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)

}

func (s *FlatCreateSuite) AfterAll(provider.T) {}

func getTokenFromBody(body io.ReadCloser) (string, error) {
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
