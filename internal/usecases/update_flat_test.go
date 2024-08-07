package usecases

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/flat"
	model_time "github.com/Ropho/avito-bootcamp-assignment/internal/models/time"
	mock_repository "github.com/Ropho/avito-bootcamp-assignment/internal/repository/mocks"
)

var (
	PosititveUpdateFlatRequest1 = FlatUpdateRequest{
		FlatID: 1,
		Status: flat.Approved,
	}
	PositiveUpdateFlat1 = flat.Model{
		FlatID:   1,
		HouseID:  1,
		Price:    1234,
		RoomsNum: 4,
		Status:   flat.Approved,
	}
)

func Test_usecase_FlatUpdate(t *testing.T) {
	type fields struct {
		repo *mock_repository.MockRepository
	}
	type args struct {
		req FlatUpdateRequest
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    FlatUpdateResponse
		wantErr error
	}{
		{
			name: "Postive #1",
			prepare: func(f *fields) {
				f.repo.EXPECT().FlatUpdate(context.Background(), PosititveUpdateFlatRequest1.FlatID, PosititveUpdateFlatRequest1.Status.String()).Return(
					PositiveUpdateFlat1, nil,
				)
			},
			args: args{
				req: PosititveUpdateFlatRequest1,
			},
			want: FlatUpdateResponse{
				FlatID:   PositiveUpdateFlat1.FlatID,
				HouseID:  PositiveUpdateFlat1.HouseID,
				Price:    PositiveUpdateFlat1.Price,
				RoomsNum: PositiveUpdateFlat1.RoomsNum,
				Status:   PositiveUpdateFlat1.Status.String(),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := &fields{
				repo: mock_repository.NewMockRepository(ctrl),
			}
			tt.prepare(f)

			u := &usecase{
				repo: f.repo,
				time: model_time.NewTimeImpl(currentTime),
			}

			got, err := u.FlatUpdate(
				context.Background(),
				tt.args.req)
			assert.ErrorIs(t, err, tt.wantErr, "errors not equal")

			assert.Equal(t, tt.want, got, "not equal response")

		})
	}
}
