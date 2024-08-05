package repository

import (
	"github.com/Ropho/avito-bootcamp-assignment/internal/models/flat"
	"github.com/Ropho/avito-bootcamp-assignment/internal/models/house"
)

type Repository interface {
	HouseCreate(house.Model) (uint32, error)
	FlatCreate(flat.Model) (uint32, error)
}
