package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/flat"
)

func (r *pgRepo) FlatUpdate(ctx context.Context, flatID uint32, flatStatus string) (flat.Model, error) {
	var err error
	var resp = flat.Model{
		FlatID: flatID,
	}
	var statusString string

	tx, err := r.conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		return flat.Model{}, fmt.Errorf("failed to init transaction: [%w]", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, updateFlatQuery, flatStatus, flatID)
	if err != nil {
		return flat.Model{}, fmt.Errorf("failed to update flat status with id %d, status %s: [%w]", flatID, flatStatus, err)
	}

	err = tx.QueryRow(getFlatByID, flatID).Scan(
		&resp.HouseID,
		&resp.Price,
		&resp.RoomsNum,
		&statusString,
	)
	if err != nil {
		return flat.Model{}, fmt.Errorf("failed to query flat with id %d: [%w]", flatID, err)
	}
	resp.Status = flat.GetStatusFromString(statusString)

	err = tx.Commit()
	if err != nil {
		return flat.Model{}, fmt.Errorf("failed to commit transaction: [%w]", err)
	}

	return resp, nil
}

var updateFlatQuery = `
UPDATE flats
SET status = $1
WHERE flatID = $2
`

var getFlatByID = `
SELECT house_id, price, rooms_number, status
FROM flats 
WHERE flat_id = $1
`
