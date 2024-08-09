package usecases

import (
	"context"
	"fmt"

	"github.com/Ropho/avito-bootcamp-assignment/pkg/sender"
)

func (u *usecases) SendEmailSubscribers(ctx context.Context, req SendEmailSubscribersRequest) error {
	var err error

	emails, err := u.repo.GetEmailsByHouseID(ctx, req.HouseID)
	if err != nil {
		return fmt.Errorf("failed to get emails to send to: [%w]", err)
	}

	message := createSubMessage()
	sender := sender.New()

	for _, email := range emails {
		err = sender.SendEmail(ctx, email, message)
		if err != nil {
			return fmt.Errorf("failed to send sub notification: [%w]", err)
		}
	}

	return nil

}

func createSubMessage() string {
	return `В доме, которым вы ранее интересовались появилась новая квартира. 
	Заходите, чтобы узнать больше!`
}
