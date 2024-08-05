package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/flat"
	model_time "github.com/Ropho/avito-bootcamp-assignment/internal/models/time"
	mock_repository "github.com/Ropho/avito-bootcamp-assignment/internal/repository/mocks"
)

var (
	PostitveGetHouseFlatsRequest1 = GetHouseFlatsRequest{
		HouseID:      1,
		OnlyApproved: true,
	}

	PostitveGetHouseModelFlatsResponse1 = []flat.Model{
		{
			FlatID:   2,
			HouseID:  1,
			Price:    1234,
			RoomsNum: 1,
			Status:   flat.Approved,
		},
	}

	PostitveGetHouseFlatsResponse1 = GetHouseFlatsResponse{
		Flats: []struct {
			FlatID  uint32
			HouseID uint32
			Price   uint32
			Rooms   uint32
			Status  string
		}{
			{
				FlatID:  2,
				HouseID: 1,
				Price:   1234,
				Rooms:   1,
				Status:  flat.ApprovedString,
			},
		},
	}

	PostitveGetHouseFlatsRequest2 = GetHouseFlatsRequest{
		HouseID:      1,
		OnlyApproved: false,
	}
	PostitveGetHouseModelFlatsResponse2 = []flat.Model{
		{
			FlatID:   1,
			HouseID:  1,
			Price:    1234,
			RoomsNum: 4,
			Status:   flat.OnModeration,
		},
		{
			FlatID:   2,
			HouseID:  1,
			Price:    1234,
			RoomsNum: 1,
			Status:   flat.Approved,
		},
	}
	PostitveGetHouseFlatsResponse2 = GetHouseFlatsResponse{
		Flats: []struct {
			FlatID  uint32
			HouseID uint32
			Price   uint32
			Rooms   uint32
			Status  string
		}{
			{
				FlatID:  2,
				HouseID: 1,
				Price:   1234,
				Rooms:   1,
				Status:  flat.ApprovedString,
			},
			{
				FlatID:  1,
				HouseID: 1,
				Price:   1234,
				Rooms:   4,
				Status:  flat.OnModeration.String(),
			},
		},
	}
)

func Test_usecase_GetHouseFlats(t *testing.T) {
	type fields struct {
		repo *mock_repository.MockRepository
	}
	type args struct {
		req GetHouseFlatsRequest
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    GetHouseFlatsResponse
		wantErr error
	}{
		{
			name: "Postive #1: regular user",
			prepare: func(f *fields) {
				f.repo.EXPECT().GetHouseFlats(PostitveGetHouseFlatsRequest1.HouseID,
					PostitveGetHouseFlatsRequest1.OnlyApproved).Return(PostitveGetHouseModelFlatsResponse1, nil)
			},
			args: args{
				req: PostitveGetHouseFlatsRequest1,
			},
			want:    PostitveGetHouseFlatsResponse1,
			wantErr: nil,
		},
		{
			name: "Postive #2: moderator",
			prepare: func(f *fields) {
				f.repo.EXPECT().GetHouseFlats(PostitveGetHouseFlatsRequest2.HouseID,
					PostitveGetHouseFlatsRequest2.OnlyApproved).Return(PostitveGetHouseModelFlatsResponse2, nil)
			},
			args: args{
				req: PostitveGetHouseFlatsRequest2,
			},
			want:    PostitveGetHouseFlatsResponse2,
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

			got, err := u.GetHouseFlats(tt.args.req)
			assert.ErrorIs(t, err, tt.wantErr, "errors not equal")

			assert.ElementsMatch(t, tt.want.Flats, got.Flats, "not equal response")

		})
	}
}
