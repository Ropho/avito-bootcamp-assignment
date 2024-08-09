package postgres

import (
	"context"
	"fmt"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/house"
)

func (r *pgRepo) HouseCreate(ctx context.Context, house house.Model) (uint32, error) {
	var err error
	var houseID int64

	err = r.conn.QueryRowContext(
		ctx,
		createHouseQuery,
		house.Address,
		house.Year,
		house.Developer,
		house.CreatedAt,
		house.UpdatedAt,
	).Scan(&houseID)
	if err != nil {
		return 0, fmt.Errorf("failed to create house %v: [%w]", house, err)
	}

	return uint32(houseID), nil
}

var createHouseQuery = `
INSERT INTO houses (address, year, developer, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING house_id
`
