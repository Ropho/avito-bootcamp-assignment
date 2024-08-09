package worker

import (
	"context"
	"fmt"

	"github.com/Ropho/avito-bootcamp-assignment/pkg/sender"
)

type worker struct {
	emailsChan chan []string
}

type NewWorkerParams struct {
	EmailsChan chan []string
}

func NewWEmailorker(p NewWorkerParams) worker {
	return worker{
		emailsChan: p.EmailsChan,
	}
}

func (w worker) Work() error {
	sender := sender.New()

	for {
		select {
		case emails := <-w.emailsChan:
			for _, email := range emails {
				err := sender.SendEmail(context.Background(), email, createSubMessage())
				if err != nil {
					return fmt.Errorf("failed to send sub notification: [%w]", err)
				}
			}

		}
	}
}

func createSubMessage() string {
	return `В доме, которым вы ранее интересовались, появилась новая квартира. 
	Заходите, чтобы узнать больше!`
}
