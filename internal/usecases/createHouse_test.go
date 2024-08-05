package usecases

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/house"
	model_time "github.com/Ropho/avito-bootcamp-assignment/internal/models/time"
	mock_repository "github.com/Ropho/avito-bootcamp-assignment/internal/repository/mocks"
)

var (
	currentTime           = time.Now()
	PostitveHouseRequest1 = HouseCreateRequest{
		Address:   "Лесная улица, 7, Москва, 125196",
		Year:      2000,
		Developer: "Мэрия города",
	}

	PostitveHouseRequest2 = HouseCreateRequest{
		Address: "Dolgoprudnny",
		Year:    2000,
	}
)

func Test_usecase_HouseCreate(t *testing.T) {
	type fields struct {
		repo *mock_repository.MockRepository
	}
	type args struct {
		req HouseCreateRequest
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    HouseCreateResponse
		wantErr error
	}{
		{
			name: "Postive #1",
			prepare: func(f *fields) {
				f.repo.EXPECT().HouseCreate(house.Model{
					HouseID:   0,
					Address:   PostitveHouseRequest1.Address,
					Year:      PostitveHouseRequest1.Year,
					Developer: PostitveHouseRequest1.Developer,
					CreatedAt: currentTime,
					UpdatedAt: currentTime,
				}).Return(
					uint32(1), nil,
				)
			},
			args: args{
				req: PostitveHouseRequest1,
			},
			want: HouseCreateResponse{
				HouseID:   1,
				Address:   PostitveHouseRequest1.Address,
				Year:      PostitveHouseRequest1.Year,
				Developer: PostitveHouseRequest1.Developer,
				CreatedAt: currentTime.String(),
				UpdatedAt: currentTime.String(),
			},
			wantErr: nil,
		},
		{
			name: "Postive #2: no developer",
			prepare: func(f *fields) {
				f.repo.EXPECT().HouseCreate(house.Model{
					HouseID:   0,
					Address:   PostitveHouseRequest2.Address,
					Year:      PostitveHouseRequest2.Year,
					Developer: PostitveHouseRequest2.Developer,
					CreatedAt: currentTime,
					UpdatedAt: currentTime,
				}).Return(
					uint32(1), nil,
				)
			},
			args: args{
				req: PostitveHouseRequest2,
			},
			want: HouseCreateResponse{
				HouseID:   1,
				Address:   PostitveHouseRequest2.Address,
				Year:      PostitveHouseRequest2.Year,
				Developer: PostitveHouseRequest2.Developer,
				CreatedAt: currentTime.String(),
				UpdatedAt: currentTime.String(),
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

			got, err := u.HouseCreate(tt.args.req)
			assert.ErrorIs(t, err, tt.wantErr, "errors not equal")

			assert.Equal(t, tt.want, got, "not equal response")

		})
	}
}
