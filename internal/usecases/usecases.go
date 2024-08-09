package usecases

import (
	"context"

	"github.com/google/uuid"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/flat"
	"github.com/Ropho/avito-bootcamp-assignment/internal/models/time"
	"github.com/Ropho/avito-bootcamp-assignment/internal/repository"
	"github.com/Ropho/avito-bootcamp-assignment/pkg/jwt"
	"github.com/Ropho/avito-bootcamp-assignment/pkg/logger"
)

type HouseCreateRequest struct {
	Address   string
	Year      uint32
	Developer string
}
type HouseCreateResponse struct {
	HouseID   uint32
	Address   string
	Year      uint32
	Developer string
	CreatedAt string
	UpdatedAt string
}

type HouseSubscribeRequest struct {
	Email   string
	UserID  uuid.UUID
	HouseID int
}
type HouseSubscribeResponse struct {
}

type FlatCreateRequest struct {
	HouseID  uint32
	Price    uint32
	RoomsNum uint32
}
type FlatCreateResponse struct {
	FlatID   uint32
	HouseID  uint32
	Price    uint32
	RoomsNum uint32
	Status   string
}

type FlatUpdateRequest struct {
	FlatID uint32
	Status flat.Status
}
type FlatUpdateResponse struct {
	FlatID   uint32
	HouseID  uint32
	Price    uint32
	RoomsNum uint32
	Status   string
}

type GetHouseFlatsRequest struct {
	HouseID      uint32
	OnlyApproved bool
}

type GetHouseFlatsResponse struct {
	Flats []struct {
		FlatID  uint32
		HouseID uint32
		Price   uint32
		Rooms   uint32
		Status  string
	}
}

type RegisterUserRequest struct {
	Email    string
	Password string
	UserType string
}
type RegisterUserResponse struct {
	UUID uuid.UUID
}

type LoginUserRequest struct {
	UUID     uuid.UUID
	Password string
}
type LoginUserResponse struct {
	Token string
}

type GetDummyLoginRequest struct {
	UserType string
}
type GetDummyLoginResponse struct {
	Token string
}

type SendEmailSubscribersRequest struct {
	HouseID int
}

type Usecases interface {
	HouseCreate(ctx context.Context, req HouseCreateRequest) (HouseCreateResponse, error)
	FlatCreate(ctx context.Context, req FlatCreateRequest) (FlatCreateResponse, error)
	FlatUpdate(ctx context.Context, req FlatUpdateRequest) (FlatUpdateResponse, error)
	GetHouseFlats(ctx context.Context, req GetHouseFlatsRequest) (GetHouseFlatsResponse, error)
	RegisterUser(ctx context.Context, req RegisterUserRequest) (RegisterUserResponse, error)
	LoginUser(ctx context.Context, req LoginUserRequest) (LoginUserResponse, error)
	GetDummyLogin(ctx context.Context, req GetDummyLoginRequest) (GetDummyLoginResponse, error)
	HouseSubscribe(ctx context.Context, req HouseSubscribeRequest) (HouseSubscribeResponse, error)
	SendEmailSubscribers(ctx context.Context, req SendEmailSubscribersRequest) error
}

type usecases struct {
	repo       repository.Repository
	time       time.Time
	jwtService jwt.Service
	logger     logger.Logger

	emailChan chan []string
}

type NewUsecasesParams struct {
	Repo       repository.Repository
	Time       time.Time
	JWTService jwt.Service
	Logger     logger.Logger
	EmailChan  chan []string
}

func NewUsecases(p NewUsecasesParams) usecases {
	return usecases{
		repo:       p.Repo,
		time:       p.Time,
		jwtService: p.JWTService,
		logger:     p.Logger,
		emailChan:  p.EmailChan,
	}
}
