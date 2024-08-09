package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/flat"
	"github.com/Ropho/avito-bootcamp-assignment/internal/models/house"
	"github.com/Ropho/avito-bootcamp-assignment/internal/models/user"
)

type Repository interface {
	HouseCreate(ctx context.Context, house house.Model) (uint32, error)
	FlatCreate(ctx context.Context, flat flat.Model) (uint32, error)
	FlatUpdate(ctx context.Context, flatID uint32, flatStatus string) (flat.Model, error)
	GetHouseFlats(ctx context.Context, houseID uint32, onlyApproved bool) ([]flat.Model, error)
	RegisterUser(ctx context.Context, user user.Model) error
	GetUserByID(ctx context.Context, uuid uuid.UUID) (user.Model, error)
}
