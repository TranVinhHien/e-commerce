package db

import (
	"context"
	"fmt"
	db "new-project/db/postgresql"
	"new-project/services"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore services.ServicesRepository

func TestMain(m *testing.M) {

	connPool, err := pgxpool.New(context.Background(), "postgresql://root:12345@localhost:5432/new_project?sslmode=disable")
	if err != nil {
		fmt.Println("db not connecting", err)
	}
	testStore = db.NewStore(connPool)
	os.Exit(m.Run())
}
