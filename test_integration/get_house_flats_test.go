//go:build !unit
// +build !unit

package testintegration

import (
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

type GetHouseFlats struct {
	suite.Suite
}

type getHouseFlatsResponse struct {
	Flats []flatResponse `json:"flats"`
}
type flatResponse struct {
	ID      int    `json:"id"`
	HouseID int    `json:"house_id"`
	Price   int    `json:"price"`
	Rooms   int    `json:"rooms"`
	Status  string `json:"status"`
}

func (s *GetHouseFlats) BeforeAll(provider.T) {}

func (s *GetHouseFlats) TestService_GetHouseFlatsModerator(t provider.T) {
	houseID := 2
	expected := []flatResponse{
		{
			ID:      5,
			HouseID: houseID,
			Price:   1,
			Rooms:   1,
			Status:  string(api.Approved),
		},
		{
			ID:      6,
			HouseID: houseID,
			Price:   2,
			Rooms:   2,
			Status:  string(api.OnModeration),
		},
		{
			ID:      7,
			HouseID: houseID,
			Price:   3,
			Rooms:   3,
			Status:  string(api.Created),
		},
		{
			ID:      8,
			HouseID: houseID,
			Price:   4,
			Rooms:   4,
			Status:  string(api.Declined),
		},
	}
	resp, err := http.Get("http://localhost:8080/dummyLogin?user_type=moderator")
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)
	defer resp.Body.Close()

	token, err := testhelpers.GetTokenFromBody(resp.Body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/house/%d", houseID), nil)
	require.NoError(t, err)
	req.Header.Add(service.AuthorizationHeader, token)

	client := http.Client{}
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)
	defer resp.Body.Close()

	got, err := getFlatsFromBody(resp.Body)
	require.NoError(t, err)

	require.ElementsMatch(t, expected, got)
}

func (s *GetHouseFlats) TestService_GetHouseFlatsClient(t provider.T) {

	houseID := 2
	expected := []flatResponse{
		{
			ID:      5,
			HouseID: houseID,
			Price:   1,
			Rooms:   1,
			Status:  string(api.Approved),
		},
	}
	resp, err := http.Get("http://localhost:8080/dummyLogin?user_type=client")
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)
	defer resp.Body.Close()

	token, err := testhelpers.GetTokenFromBody(resp.Body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/house/%d", houseID), nil)
	require.NoError(t, err)
	req.Header.Add(service.AuthorizationHeader, token)

	client := http.Client{}
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)
	defer resp.Body.Close()

	got, err := getFlatsFromBody(resp.Body)
	require.NoError(t, err)

	require.ElementsMatch(t, expected, got)

}
func (s *GetHouseFlats) AfterAll(provider.T) {}

func getFlatsFromBody(body io.ReadCloser) ([]flatResponse, error) {

	var resp getHouseFlatsResponse
	bytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: [%w]", err)
	}
	defer body.Close()

	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json struct: [%w]", err)
	}

	return resp.Flats, nil
}
