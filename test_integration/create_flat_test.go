//go:build !unit
// +build !unit

package testintegration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Ropho/avito-bootcamp-assignment/api"
	"github.com/Ropho/avito-bootcamp-assignment/internal/service"
	"github.com/Ropho/avito-bootcamp-assignment/test_integration/testhelpers"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/require"
)

type FlatCreateSuite struct {
	suite.Suite
}

type createFlatResponse struct {
	ID      int    `json:"id"`
	HouseID int    `json:"house_id"`
	Price   int    `json:"price"`
	Rooms   int    `json:"rooms"`
	Status  string `json:"status"`
}

func (s *FlatCreateSuite) BeforeAll(provider.T) {}

func (s *FlatCreateSuite) TestService_FlatCreateModerator(t provider.T) {

	type bodyFlatCreateType struct {
		HouseID int32 `json:"house_id"`
		Price   int32 `json:"price"`
		Rooms   int32 `json:"rooms"`
	}

	resp, err := http.Get("http://localhost:8080/dummyLogin?user_type=moderator")
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)
	defer resp.Body.Close()

	token, err := testhelpers.GetTokenFromBody(resp.Body)
	require.NoError(t, err)

	body := bodyFlatCreateType{
		HouseID: 1,
		Price:   1000,
		Rooms:   12,
	}

	bodyBytes, err := json.Marshal(&body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/flat/create", bytes.NewReader(bodyBytes))
	require.NoError(t, err)
	req.Header.Set(service.AuthorizationHeader, token)

	client := http.Client{}
	resp, err = client.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()
	require.Equal(t, resp.StatusCode, http.StatusOK)

	got, err := getFlatFromBody(resp.Body)
	require.NoError(t, err)

	expected := createFlatResponse{
		HouseID: int(body.HouseID),
		Price:   int(body.Price),
		Rooms:   int(body.Rooms),
		Status:  string(api.Created),
	}

	require.Equal(t, expected.HouseID, got.HouseID)
	require.Equal(t, expected.Price, got.Price)
	require.Equal(t, expected.Rooms, got.Rooms)
	require.Equal(t, expected.Status, got.Status)

}

func (s *FlatCreateSuite) TestService_FlatCreateClient(t provider.T) {

	type bodyFlatCreateType struct {
		HouseID int32 `json:"house_id"`
		Price   int32 `json:"price"`
		Rooms   int32 `json:"rooms"`
	}

	resp, err := http.Get("http://localhost:8080/dummyLogin?user_type=client")
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)
	defer resp.Body.Close()

	token, err := testhelpers.GetTokenFromBody(resp.Body)
	require.NoError(t, err)

	body := bodyFlatCreateType{
		HouseID: 1,
		Price:   1234,
		Rooms:   1,
	}

	bodyBytes, err := json.Marshal(&body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/flat/create", bytes.NewReader(bodyBytes))
	require.NoError(t, err)
	req.Header.Set(service.AuthorizationHeader, token)

	client := http.Client{}
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, resp.StatusCode, http.StatusOK)

	got, err := getFlatFromBody(resp.Body)
	require.NoError(t, err)

	expected := createFlatResponse{
		HouseID: int(body.HouseID),
		Price:   int(body.Price),
		Rooms:   int(body.Rooms),
		Status:  string(api.Created),
	}

	require.Equal(t, expected.HouseID, got.HouseID)
	require.Equal(t, expected.Price, got.Price)
	require.Equal(t, expected.Rooms, got.Rooms)
	require.Equal(t, expected.Status, got.Status)
}
func (s *FlatCreateSuite) AfterAll(provider.T) {}

func getFlatFromBody(body io.ReadCloser) (createFlatResponse, error) {

	var resp createFlatResponse
	bytes, err := io.ReadAll(body)
	if err != nil {
		return createFlatResponse{}, fmt.Errorf("failed to read body: [%w]", err)
	}

	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return createFlatResponse{}, fmt.Errorf("failed to unmarshal json struct: [%w]", err)
	}

	return resp, nil
}
