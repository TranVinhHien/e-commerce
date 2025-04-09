package db

import (
	"context"
	"database/sql"
	"fmt"
	db "new-project/db/sqlc"
	"new-project/services"
)

type SQLStore struct {
	*db.Queries
	connPool *sql.DB
}

type Store interface {
	db.Querier
}

// create new store

func NewStore(connectPool *sql.DB) services.ServicesRepository {
	return &SQLStore{
		Queries:  db.New(connectPool),
		connPool: connectPool,
	}
}

// write a function transaction using package github.com/jackc/pgx/v5/pgxpool
func (s *SQLStore) execTS(ctx context.Context, fn func(tx *db.Queries) error) error {
	tx, err := s.connPool.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := db.New(tx)
	err = fn(q)
	if err != nil {
		fmt.Printf("lôi err")
		if errTran := tx.Rollback(); errTran != nil {
			return fmt.Errorf("transaction error %v ,rollback trancsaction error : %v", err, errTran)
		}
		return err
	}

	return tx.Commit()
}
