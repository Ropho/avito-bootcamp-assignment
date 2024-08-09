package router

import (
	"context"
	"net/http"
	"strings"

	"github.com/Ropho/avito-bootcamp-assignment/api"
	"github.com/Ropho/avito-bootcamp-assignment/internal/service"
	"github.com/Ropho/avito-bootcamp-assignment/pkg/jwt"
	"github.com/Ropho/avito-bootcamp-assignment/pkg/logger"
	"github.com/google/uuid"
)

const (
	getHouseFlats  = "house/"
	flatUpdate     = "flat/update"
	houseCreate    = "house/create"
	flatCreate     = "flat/create"
	houseSubscribe = "/subscribe"
)

func (m *Manager) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context

		// очень банальная обработка
		needsAutentication :=
			strings.Contains(r.URL.Path, flatUpdate) ||
				strings.Contains(r.URL.Path, houseCreate) ||
				strings.Contains(r.URL.Path, flatCreate) ||
				strings.Contains(r.URL.Path, getHouseFlats) ||
				strings.Contains(r.URL.Path, houseSubscribe)

		if needsAutentication {

			userClaims, err := m.jwtService.VerifyUserAccessToken(r.Header.Get(service.AuthorizationHeader))
			if err != nil {
				m.logger.Errorf(err, "failed to parse authorization token")

				err = service.WriteUnauthorized(w)
				if err != nil {
					m.logger.Fatal("failed to send response", err)
				}
				return
			}

			ctx = context.WithValue(r.Context(), service.UserIDKey{}, userClaims.UUID)
			ctx = context.WithValue(ctx, service.UserTypeKey{}, userClaims.UserType)
			r = r.WithContext(ctx)

			if strings.Contains(r.URL.Path, flatUpdate) ||
				strings.Contains(r.URL.Path, houseCreate) {
				if userClaims.UserType != string(api.Moderator) {
					m.logger.Errorf(err, "moderator required")

					err = service.WriteUnauthorized(w)
					if err != nil {
						m.logger.Fatal("failed to send response", err)
					}
					return
				}
			}

		}

		ctx = context.WithValue(r.Context(), service.RequestIDKey{}, m.generateRequestID())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Manager) generateRequestID() string {
	uuid, err := uuid.NewV6()
	if err != nil {
		m.logger.Fatal("failed to generate request id: %v", err)
	}

	return uuid.String()
}

type Manager struct {
	logger     logger.Logger
	jwtService jwt.Service
}

func NewInterceptorsManager(jwtService jwt.Service, logger logger.Logger) Manager {
	return Manager{jwtService: jwtService, logger: logger}
}
