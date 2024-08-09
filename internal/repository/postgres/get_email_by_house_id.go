package postgres

import (
	"context"
	"fmt"
)

func (r *pgRepo) GetEmailsByHouseID(ctx context.Context, houseID int) ([]string, error) {
	var err error
	var emails []string
	var curEmail string

	rows, err := r.conn.QueryContext(
		ctx,
		getEmailsByHouseID,
		houseID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get emails by house id %v: [%w]", houseID, err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&curEmail,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan email: [%w]", err)
		}

		emails = append(emails, curEmail)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to scan from rows: [%w]", err)
	}

	return emails, nil
}

var getEmailsByHouseID = `
SELECT email FROM subscriptions
WHERE house_id = $1
`
