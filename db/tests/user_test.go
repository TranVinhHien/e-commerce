package db

// import (
// 	"context"
// 	"fmt"
// 	util_assets "new-project/assets/util"
// 	db "new-project/db/sqlc"
// 	"testing"
// 	"time"

// 	"github.com/jackc/pgx/v5/pgtype"
// 	"github.com/stretchr/testify/require"
// )

// func TestAddUser(t *testing.T) {

// 	user, err := testStore.GetUser(context.Background(), "r5ZxKonFoF")
// 	require.NoError(t, err)
// 	fmt.Print(user)
// 	require.Empty(t, user)
// }

// // func Test Create User
// func TestCreateUser(t *testing.T) {

// 	userName123 := util_assets.RandomString(10)
// 	user, err := testStore.(context.Background(),
// 		db.CreateUserParams{
// 			Username: userName123,
// 			Password: util_assets.RandomString(10),
// 			FullName: util_assets.RandomString(20),
// 			CreateAt: time.Now(),
// 		})
// 	require.NoError(t, err)
// 	require.Equal(t, user.Username, userName123)
// }
// // func TestUpdateUser(t *testing.T) {
// // 	userName123 := "r5ZxKonFoF"
// // 	Password := util_assets.RandomString(20)
// // 	user, err := testStore.UpdateUser(context.Background(),
// // 		db.UpdateUserParams{
// // 			Username: userName123,
// // 			Password: pgtype.Text{String: Password, Valid: true},
// // 			// FullName: pgtype.Text{String: FullName, Valid: true},

// // 			IsActive: pgtype.Bool{Valid: true, Bool: true},
// // 		})
// // 	require.NoError(t, err)
// // 	require.Equal(t, userName123, user.Username)
// // 	require.Equal(t, "yFfhpInNHxuu97Bd9aCI", user.FullName)

// // }
