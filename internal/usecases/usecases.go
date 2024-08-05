package usecases

import (
	"github.com/Ropho/avito-bootcamp-assignment/internal/models/flat"
	"github.com/Ropho/avito-bootcamp-assignment/internal/models/time"
	"github.com/Ropho/avito-bootcamp-assignment/internal/repository"
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

type Usecase interface {
	HouseCreate(HouseCreateRequest) (HouseCreateResponse, error)
	FlatCreate(FlatCreateRequest) (FlatCreateResponse, error)
	FlatUpdate(FlatUpdateRequest) (FlatUpdateResponse, error)
}

type usecase struct {
	repo repository.Repository
	time time.Time
}

type NewUsecaseParams struct {
	Repo repository.Repository
	Time time.Time
}

func NewUsecase(p NewUsecaseParams) usecase {
	return usecase{
		repo: p.Repo,
		time: p.Time,
	}
}
