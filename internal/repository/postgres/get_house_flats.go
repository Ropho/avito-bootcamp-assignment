package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/flat"
)

func (r *pgRepo) GetHouseFlats(ctx context.Context, houseID uint32, onlyApproved bool) ([]flat.Model, error) {
	var err error
	var rows *sql.Rows
	var resp []flat.Model

	if onlyApproved {
		rows, err = r.conn.QueryContext(
			ctx,
			getHouseFlatsApprovedOnlyQuery,
			houseID,
			flat.Approved.String(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to query get house flats approved only: [%w]", err)
		}
	} else {
		rows, err = r.conn.QueryContext(
			ctx,
			getHouseFlatsQuery,
			houseID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to query get house flats: [%w]", err)
		}
	}

	for rows.Next() {
		var statusString string
		curFlat := flat.Model{
			HouseID: houseID,
		}

		err = rows.Scan(
			&curFlat.FlatID,
			&curFlat.Price,
			&curFlat.RoomsNum,
			&statusString,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: [%w]", err)
		}
		curFlat.Status = flat.GetStatusFromString(statusString)

		resp = append(resp, curFlat)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate rows: [%w]", err)
	}

	return resp, nil
}

var getHouseFlatsApprovedOnlyQuery = `
SELECT flat_id, price, rooms_number, status FROM flats
WHERE house_id = $1 AND status = $2
`
var getHouseFlatsQuery = `
SELECT flat_id, price, rooms_number, status FROM flats
WHERE house_id = $1
`
