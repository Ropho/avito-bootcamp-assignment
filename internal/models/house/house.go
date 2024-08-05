package house

import (
	"time"

	ownTime "github.com/Ropho/avito-bootcamp-assignment/internal/models/time"
)

type Model struct {
	HouseID   uint32
	Address   string
	Year      uint32
	Developer string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewParams struct {
	Address   string
	Year      uint32
	Developer string
	Time      ownTime.Time
}

func New(params NewParams) Model {
	createdAt := params.Time.Now()

	house := Model{
		Address:   params.Address,
		Year:      params.Year,
		Developer: params.Developer,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	return house
}
