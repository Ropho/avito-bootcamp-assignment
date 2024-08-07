package postgres_test

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ropho/avito-bootcamp-assignment/internal/models/flat"
	"github.com/Ropho/avito-bootcamp-assignment/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var flatsResponse1 = []flat.Model{
	{
		FlatID:   1,
		HouseID:  1,
		Price:    1234,
		RoomsNum: 4,
		Status:   flat.Approved,
	},
	{
		FlatID:   2,
		HouseID:  1,
		Price:    12345,
		RoomsNum: 3,
		Status:   flat.OnModeration,
	},
	{
		FlatID:   3,
		HouseID:  1,
		Price:    123411,
		RoomsNum: 3,
		Status:   flat.Created,
	},
}

func TestGetHouseFlats(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "failed to init sql db mock: ", err)
	defer db.Close()

	repo := postgres.NewPgRepo(
		&postgres.NewPgRepoParams{
			Conn: db,
		},
	)

	type args struct {
		houseID      uint32
		onlyApproved bool
	}

	type mockBehaviour func(args args, wantFlats []flat.Model)

	testTable := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		wantFlats     []flat.Model
		wantErr       error
	}{
		{
			name: "Positive #1: admin request",
			args: args{
				houseID:      1,
				onlyApproved: false,
			},
			mockBehaviour: func(args args, wantFlats []flat.Model) {

				var expectedValues [][]driver.Value
				for _, flat := range wantFlats {
					expectedValues = append(expectedValues,
						[]driver.Value{flat.FlatID, flat.Price, flat.RoomsNum, flat.Status.String()})
				}

				rows := sqlmock.NewRows([]string{"flat_id", "price", "rooms_number", "status"}).
					AddRows(expectedValues...)

				mock.ExpectQuery(`SELECT FROM flats`).WithArgs(
					args.houseID,
				).WillReturnRows(rows)
			},
			wantFlats: flatsResponse1,
			wantErr:   nil,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			t.Helper()
			test.mockBehaviour(test.args, test.wantFlats)

			gotFlats, err := repo.GetHouseFlats(context.Background(), test.args.houseID, test.args.onlyApproved)

			require.ElementsMatch(t, test.wantFlats, gotFlats, "flats not equal")
			require.ErrorIs(t, err, test.wantErr, "errors not equal ", err, " ", test.wantErr)
		})
	}
}
