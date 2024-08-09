package postgres

import (
	"context"
	"fmt"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/subscription"
)

func (r *pgRepo) SubscriptionCreate(ctx context.Context, sub subscription.Model) error {
	var err error

	_, err = r.conn.ExecContext(
		ctx,
		createSubsciptionQuery,
		sub.UserID,
		sub.HouseID,
		sub.Email,
	)
	if err != nil {
		return fmt.Errorf("failed to create subscription %v: [%w]", sub, err)
	}

	return nil
}

var createSubsciptionQuery = `
INSERT INTO subscriptions (user_id, house_id, email)
VALUES ($1, $2, $3)
`
