package router

import (
	"context"
)

type UserAuthorizedContext struct {
	context.Context
	UUID      string
	RequestID string
}
