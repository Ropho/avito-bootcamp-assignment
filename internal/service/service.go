package service

import (
	"github.com/Ropho/avito-bootcamp-assignment/internal/usecases"
	"github.com/Ropho/avito-bootcamp-assignment/pkg/logger"
)

const AuthorizationHeader = "Authorization"

type UserIDKey struct{}
type UserTypeKey struct{}
type RequestIDKey struct{}

type Service struct {
	usecases usecases.Usecases
	logger   logger.Logger
}

type NewServiceParams struct {
	Usecases usecases.Usecases
	Logger   logger.Logger
}

func NewService(p NewServiceParams) Service {
	return Service{
		usecases: p.Usecases,
		logger:   p.Logger,
	}
}
