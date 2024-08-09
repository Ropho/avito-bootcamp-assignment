package postgres

import (
	"context"
	"fmt"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/flat"
)

func (r *pgRepo) FlatCreate(ctx context.Context, flat flat.Model) (uint32, error) {
	var err error
	var flatID int64

	err = r.conn.QueryRowContext(
		ctx,
		createFlatQuery,
		flat.HouseID,
		flat.Price,
		flat.RoomsNum,
		flat.Status.String(),
	).Scan(
		&flatID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create flat %v: [%w]", flat, err)
	}

	return uint32(flatID), nil
}

var createFlatQuery = `
INSERT INTO flats (house_id, price, rooms_number, status)
VALUES ($1, $2, $3, $4)
RETURNING flat_id
`
