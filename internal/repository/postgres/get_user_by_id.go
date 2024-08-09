package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/user"
)

func (r *pgRepo) GetUserByID(ctx context.Context, uuid uuid.UUID) (user.Model, error) {
	var err error
	var userTypeString string

	userModel := user.Model{
		ID: uuid,
	}

	err = r.conn.QueryRowContext(
		ctx,
		getUserByIDQuery,
		userModel.ID,
	).Scan(
		&userModel.Email,
		&userModel.EncryptedPassword,
		&userModel.Salt,
		&userTypeString,
	)
	if err != nil {
		return user.Model{}, fmt.Errorf("failed to get user by id %v: [%w]", userModel, err)
	}

	userModel.Type = user.GetTypeFromString(userTypeString)

	return userModel, nil
}

var getUserByIDQuery = `
SELECT email, encr_pass, salt, user_type FROM users
WHERE "uuid" = $1
`
