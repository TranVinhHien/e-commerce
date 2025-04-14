package db

import (
	"context"
	"database/sql"
	"fmt"
	db "new-project/db/sqlc"
	"new-project/services"
	entity "new-project/services/entity"
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
func listData(ctx context.Context, connPool *sql.DB, table string, query entity.QueryFilter) (*sql.Rows, error) {
	querySQL := fmt.Sprintf("SELECT *  FROM %s WHERE 1=1", table)
	args := []interface{}{}

	// Xây dựng SQL từ QueryFilter
	for _, condition := range query.Conditions {
		querySQL += fmt.Sprintf(" AND %s %s ?", condition.Field, condition.Operator)
		args = append(args, condition.Value)
	}

	// Thêm sắp xếp
	if query.OrderBy != nil {
		querySQL += fmt.Sprintf(" ORDER BY %s %s", query.OrderBy.Field, query.OrderBy.Value)
	}

	// Thêm phân trang
	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		querySQL += " LIMIT ? OFFSET ?"
		args = append(args, query.PageSize, offset)
	}

	// Thực thi truy vấn

	rows, err := connPool.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
