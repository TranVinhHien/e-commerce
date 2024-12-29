package main

// import gin
import (
	"context"
	config_assets "new-project/assets/config"
	"new-project/assets/token"
	"new-project/controllers"
	db "new-project/db/postgresql"
	redis_db "new-project/db/redis"
	"new-project/services"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// create logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// read config file
	env, err := config_assets.LoadConfig(".")
	if err != nil {
		log.Err(err).Msg("Error read env:")
		return
	}
	// create connect to database
	conn, err := connectDBWithRetry(5, env.DBSource)
	if conn == nil {
		log.Err(err).Msg("Error when created connect to database")
		return
	}
	// close connection after gin stopped
	defer conn.Close()

	log.Info().Msg("Connect to database successfully")
	db := db.NewStore(conn)
	log.Info().Msg("Creating gin server...")
	// create jwt
	jwtMaker, err := token.NewJWTMaker(env.JWTSecret)
	if err != nil {
		log.Err(err).Msg("Error create JWTMaker")
		return
	}
	rdb, err := connectDBRedisWithRetry(5)
	if err != nil {
		log.Err(err).Msg("Error when created connect to redis")
		return
	}
	//setup redis Options
	redisdb := redis_db.NewRedisDB(rdb)
	// start jobs
	go redisdb.RemoveTokenExp(redis_db.BLACK_LIST)

	services := services.NewService(db, jwtMaker, env, redisdb)
	controller := controllers.NewAPIController(services, jwtMaker)

	engine := gin.Default()
	// config cors middleware
	config := cors.Config{
		AllowOrigins:     []string{env.ClientIP},                              // Chỉ cho phép localhost:3000
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Các method được phép
		AllowHeaders:     []string{"Content-Type", "Origin", "Authorization"}, // Các headers được phép
		ExposeHeaders:    []string{"Content-Length"},                          // Các headers trả về
		AllowCredentials: true,                                                // Cho phép cookies
	}

	v1 := engine.Group("/v1")
	v1.Use(cors.New(config))
	controller.SetUpRoute(v1)

	log.Info().Msg("Starting server on port " + env.HTTPServerAddress)

	engine.Run(env.HTTPServerAddress)

}
func connectDBRedisWithRetry(times int) (*redis.Client, error) {
	var e error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2*time.Duration(times))
	defer cancel()
	for i := 1; i <= times; i++ {
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		_, err := rdb.Ping(ctx).Result()

		if err != nil {
			log.Err(err).Msg("Can't connect to redis")
		}
		// defer conn.Release()

		if err == nil {
			return rdb, nil
		}
		e = err
		time.Sleep(time.Second * 2)
	}
	return nil, e

}
func connectDBWithRetry(times int, dbConfig string) (*pgxpool.Pool, error) {
	var e error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2*time.Duration(times))
	defer cancel()
	for i := 1; i <= times; i++ {
		pool, err := pgxpool.New(context.Background(), dbConfig)
		if err != nil {
			log.Err(err).Msg("Can't create database pool")
		}
		_, err = pool.Acquire(ctx)
		if err != nil {
			log.Err(err).Msg("Can't get connection to database pool")
		}
		// defer conn.Release()

		if err == nil {
			return pool, nil
		}
		e = err
		time.Sleep(time.Second * 2)
	}
	return nil, e
}
