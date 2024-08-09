package postgres

import (
	"context"
	"fmt"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/user"
)

func (r *pgRepo) RegisterUser(ctx context.Context, user user.Model) error {
	var err error

	_, err = r.conn.ExecContext(
		ctx,
		createUserQuery,
		user.ID,
		user.Email,
		user.EncryptedPassword,
		user.Salt,
		user.Type.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to create user %v: [%w]", user, err)
	}

	return nil
}

var createUserQuery = `
INSERT INTO users ("uuid", email, encr_pass, salt, user_type)
VALUES ($1, $2, $3, $4, $5)
`
