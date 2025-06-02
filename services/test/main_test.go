package services_test

import (
	"context"
	"database/sql"
	config_assets "new-project/assets/config"
	assets_firebase "new-project/assets/fire-base"
	assets_jobs "new-project/assets/jobs"
	"new-project/assets/token"
	db "new-project/db/mysql"
	redis_db "new-project/db/redis"
	"new-project/services"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

var testService services.ServiceUseCase

func TestMain(m *testing.M) {

	pool, _ := sql.Open("mysql", "root:12345@tcp(localhost:3306)/e-commerce?parseTime=true")

	env, _ := config_assets.LoadConfig("../../")
	db := db.NewStore(pool)
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rerids := redis_db.NewRedisDB(rdb)
	jwtMaker, _ := token.NewJWTMaker(env.JWTSecret)
	// firebaseApp, _ := firebase.N(context.Background(), nil, )
	firebase, _ := assets_firebase.NewFirebase(context.Background(), "../../assets/firebase/credentials.json")
	job, _ := assets_jobs.NewJobScheduler()
	testService = services.NewService(db, jwtMaker, env, rerids, firebase, job)

	os.Exit(m.Run())
}
