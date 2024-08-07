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

func TestCreateFlat(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "failed to init sql db mock: ", err)
	defer db.Close()

	repo := postgres.NewPgRepo(
		&postgres.NewPgRepoParams{
			Conn: db,
		},
	)

	type args struct {
		flat flat.Model
	}

	type mockBehaviour func(args args, wantID uint32)

	testTable := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		wantID        uint32
		wantErr       error
	}{
		{
			name: "Positive #1",
			args: args{
				flat.Model{
					HouseID:  1,
					Price:    1234,
					RoomsNum: 4,
					Status:   flat.Created,
				},
			},
			mockBehaviour: func(args args, wantID uint32) {

				rows := sqlmock.NewRows([]string{"flat_id"}).AddRow(wantID)
				mock.ExpectQuery(`INSERT INTO flats`).WithArgs(
					args.flat.HouseID,
					args.flat.Price,
					args.flat.RoomsNum,
					args.flat.Status.String()).WillReturnRows(rows)
			},
			wantID:  1,
			wantErr: nil,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			t.Helper()
			test.mockBehaviour(test.args, test.wantID)

			gotID, err := repo.CreateFlat(context.Background(), test.args.flat)

			require.Equal(t, test.wantID, gotID, "ids not equal")
			require.ErrorIs(t, err, test.wantErr, "errors not equal ", err, " ", test.wantErr)
		})
	}
}
