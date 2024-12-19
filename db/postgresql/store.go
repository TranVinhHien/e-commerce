package db

import (
	db "new-project/db/sqlc"
	"new-project/services"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLStore struct {
	*db.Queries
	connPool *pgxpool.Pool
}

type Store interface {
	db.Querier
}

func NewStore(connectPool *pgxpool.Pool) services.ServicesRepository {
	return &SQLStore{
		Queries:  db.New(connectPool),
		connPool: connectPool,
	}
}

// write a function transaction using package github.com/jackc/pgx/v5/pgxpool
// func (s *SQLStore) execTS(ctx context.Context, fn func(tx *Queries) error) error {
// 	tx, err := s.connPool.Begin(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	q := New(tx)
// 	err = fn(q)
// 	if err != nil {
// 		if errTran := tx.Rollback(ctx); errTran != nil {
// 			return fmt.Errorf("transaction error %v ,rollback trancsaction error : %v", err, errTran)
// 		}
// 		return err
// 	}

// 	return tx.Commit(ctx)
// }
