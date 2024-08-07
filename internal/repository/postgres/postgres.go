package postgres

import (
	"database/sql"
	"fmt"

	// pg driver export
	_ "github.com/jackc/pgx/v5/stdlib"
)

type pgRepo struct {
	conn *sql.DB
}

// NewPgRepoParams defines necessary params to initialize posgres repository
type NewPgRepoParams struct {
	Conn *sql.DB
}

// NewPgRepo initializes Postgres Repository
func NewPgRepo(p *NewPgRepoParams) pgRepo {
	return pgRepo{
		conn: p.Conn,
	}
}

func OpenConnection(dbURL string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to init connection to db: [%w]", err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: [%w]", err)
	}

	return conn, nil
}
