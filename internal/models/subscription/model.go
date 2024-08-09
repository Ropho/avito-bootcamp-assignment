package subscription

import "github.com/google/uuid"

type Model struct {
	UserID  uuid.UUID
	HouseID int
	Email   string
}

type NewParams struct {
	UserID  uuid.UUID
	HouseID int
	Email   string
}

func New(params NewParams) Model {
	return Model(params)
}
