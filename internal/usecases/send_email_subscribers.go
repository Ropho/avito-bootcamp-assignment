package usecases

import (
	"context"
	"fmt"
)

func (u *usecases) SendEmailSubscribers(ctx context.Context, req SendEmailSubscribersRequest) error {
	var err error

	emails, err := u.repo.GetEmailsByHouseID(ctx, req.HouseID)
	if err != nil {
		return fmt.Errorf("failed to get emails to send to: [%w]", err)
	}

	if emails != nil {
		u.emailChan <- emails
	}

	return nil

}
