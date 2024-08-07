package postgres_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ropho/avito-bootcamp-assignment/internal/models/flat"
	"github.com/Ropho/avito-bootcamp-assignment/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testFlat1 = flat.Model{
	FlatID:   1,
	HouseID:  1,
	Price:    1111,
	RoomsNum: 1,
	Status:   flat.OnModeration,
}

func TestUpdateFlat(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "failed to init sql db mock: ", err)
	defer db.Close()

	repo := postgres.NewPgRepo(
		&postgres.NewPgRepoParams{
			Conn: db,
		},
	)

	type args struct {
		flatID     uint32
		flatStatus string
	}

	type mockBehaviour func(args args, want flat.Model)

	testTable := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		want          flat.Model
		wantErr       error
	}{
		{
			name: "Positive #1",
			args: args{
				flatID:     1,
				flatStatus: testFlat1.Status.String(),
			},
			mockBehaviour: func(args args, want flat.Model) {

				mock.ExpectBegin()

				mock.ExpectExec("UPDATE flats SET status").WithArgs(
					args.flatStatus,
					args.flatID,
				).WillReturnResult(sqlmock.NewErrorResult(nil))

				rows := sqlmock.NewRows([]string{"house_id", "price", "rooms_number", "status"}).
					AddRow(want.HouseID, want.Price, want.RoomsNum, want.Status.String())
				mock.ExpectQuery(`SELECT house_id, price, rooms_number, status FROM flats`).WithArgs(
					args.flatID,
				).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			want:    testFlat1,
			wantErr: nil,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			t.Helper()
			test.mockBehaviour(test.args, test.want)

			got, err := repo.FlatUpdate(context.Background(), test.args.flatID, test.args.flatStatus)

			require.ErrorIs(t, err, test.wantErr, "errors not equal ", err, " ", test.wantErr)
			require.Equal(t, test.want, got, "flats not equal")
		})
	}
}
