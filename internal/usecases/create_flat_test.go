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
	PosititveFlatRequest1 = FlatCreateRequest{
		HouseID:  1,
		Price:    1234,
		RoomsNum: 4,
	}
)

func Test_usecase_FlatCreate(t *testing.T) {
	type fields struct {
		repo *mock_repository.MockRepository
	}
	type args struct {
		req FlatCreateRequest
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    FlatCreateResponse
		wantErr error
	}{
		{
			name: "Postive #1",
			prepare: func(f *fields) {
				f.repo.EXPECT().FlatCreate(context.Background(), flat.Model{
					HouseID:  PosititveFlatRequest1.HouseID,
					Price:    PosititveFlatRequest1.Price,
					RoomsNum: PosititveFlatRequest1.RoomsNum,
					Status:   flat.Created,
				}).Return(
					uint32(1), nil,
				)
			},
			args: args{
				req: PosititveFlatRequest1,
			},
			want: FlatCreateResponse{
				FlatID:   1,
				HouseID:  PosititveFlatRequest1.HouseID,
				Price:    PosititveFlatRequest1.Price,
				RoomsNum: PosititveFlatRequest1.RoomsNum,
				Status:   flat.Created.String(),
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

			got, err := u.FlatCreate(context.Background(), tt.args.req)
			assert.ErrorIs(t, err, tt.wantErr, "errors not equal")

			assert.Equal(t, tt.want, got, "not equal response")

		})
	}
}
