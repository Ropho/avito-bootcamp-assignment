package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ropho/avito-bootcamp-assignment/internal/models/house"
	"github.com/Ropho/avito-bootcamp-assignment/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHouseCreate(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "failed to init sql db mock: ", err)
	defer db.Close()

	repo := postgres.NewPgRepo(
		&postgres.NewPgRepoParams{
			Conn: db,
		},
	)

	type args struct {
		house house.Model
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
				house.Model{
					Address:   "Moscow",
					Year:      2024,
					Developer: "",
					CreatedAt: time.Unix(123, 123),
					UpdatedAt: time.Unix(123, 123),
				},
			},
			mockBehaviour: func(args args, wantID uint32) {

				rows := sqlmock.NewRows([]string{"house_id"}).AddRow(wantID)
				mock.ExpectQuery(`INSERT INTO houses`).WithArgs(
					args.house.Address,
					args.house.Year,
					args.house.Developer,
					args.house.CreatedAt,
					args.house.UpdatedAt,
				).WillReturnRows(rows)
			},
			wantID:  1,
			wantErr: nil,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			t.Helper()
			test.mockBehaviour(test.args, test.wantID)

			gotID, err := repo.HouseCreate(context.Background(), test.args.house)

			require.ErrorIs(t, err, test.wantErr, "errors not equal ", err, " ", test.wantErr)
			require.Equal(t, test.wantID, gotID, "ids not equal")
		})
	}
}
