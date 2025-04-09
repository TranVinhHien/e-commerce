package db_test

// import (
// 	"context"
// 	"database/sql"
// 	db "new-project/db/sqlc"
// 	"testing"
// 	"time"

// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/stretchr/testify/require"
// )

// func createRandomUser() db.Users {
// 	return db.Users{
// 		Username: "testuser",
// 		Password: "password",
// 		FullName: "Test User",
// 		IsActive: true,
// 		CreateAt: time.Now(),
// 	}
// }

// func TestGetUser(t *testing.T) {
// 	pool, err := sql.Open("mysql", "root:12345@tcp(localhost:3306)/e-commerce?parseTime=true")
// 	store := db.New(pool)

// 	// Create a random user
// 	// user := createRandomUser()

// 	// Insert the user into the database

// 	require.NoError(t, err)

// 	// Retrieve the user from the database
// 	err = store.UpdateUser(context.Background(), db.UpdateUserParams{
// 		FullName: sql.NullString{
// 			String: "DMM",
// 			Valid:  true,
// 		},
// 		Username: "hien123",
// 	})
// 	require.NoError(t, err)

// 	// // Verify the retrieved user's information
// 	// require.Equal(t, "hien123", retrievedUser.Username)
// 	// require.Equal(t, user.Password, retrievedUser.Password)
// 	// require.Equal(t, user.FullName, retrievedUser.FullName)
// 	// require.Equal(t, user.IsActive, retrievedUser.IsActive)
// 	// require.WithinDuration(t, user.CreateAt, retrievedUser.CreateAt, time.Second)
// }
