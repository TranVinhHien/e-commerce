package redis_db

import (
	"context"
	"fmt"
	"new-project/services"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisDB struct {
	client *redis.Client
}
type RedisDBClient interface {
	// Save
}

func NewRedisDB(rdb *redis.Client) services.ServicesRedis {
	return &RedisDB{client: rdb}
}

// Hàm thêm token vào ZSET với thời gian hết hạn
func (s *RedisDB) addScoreMember(ctx context.Context, zsetKey, token string, expiry float64) error {
	_, err := s.client.ZAdd(ctx, zsetKey, redis.Z{
		Score: expiry,
		//Score:  float64(time.Now().Add(time.Second * 15).Unix()),
		Member: token,
	}).Result()
	if err != nil {
		return fmt.Errorf("error when add new item %s: to key storeMember:%s ,err: %v", token, zsetKey, err)
	}
	return nil
}

// scan item expired and remove it
func (s *RedisDB) removeExpired(ctx context.Context, zsetKey string) error {
	now := float64(time.Now().Unix())
	// Lấy danh sách token hết hạn
	expiredTokens, err := s.client.ZRangeByScore(ctx, zsetKey, &redis.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprintf("%f", now),
	}).Result()
	if err != nil {
		return fmt.Errorf("error when check token: %v", err)
	}
	if len(expiredTokens) == 0 {
		return nil
	}
	// Xóa các token hết hạn
	_, err = s.client.ZRem(ctx, zsetKey, expiredTokens).Result()
	if err != nil {

		return fmt.Errorf("error when remove token: %v", err)
	}
	log.Info().Msg(fmt.Sprintf("remove %v token", len(expiredTokens)))
	return nil
}

// check token Valid
func (s *RedisDB) isExists(ctx context.Context, zsetKey, token string) bool {
	// Dùng ZSCORE để kiểm tra token
	now := float64(time.Now().Unix()) // Lấy thời gian hiện tại

	score, err := s.client.ZScore(ctx, zsetKey, token).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		return false
	}
	if now >= score {
		return false
	}
	// Token tồn tại
	return true
}
func (s *RedisDB) AddTokenToBlackList(ctx context.Context, token string, exprid float64) error {
	return s.addScoreMember(ctx, BLACK_LIST, token, exprid)
}
func (s *RedisDB) CheckExistsFromBlackList(ctx context.Context, token string, exprid float64) bool {
	return s.isExists(ctx, BLACK_LIST, token)
}

// auto remove token expired
func (s *RedisDB) RemoveTokenExp(zsetKey string) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		fmt.Print(err)
	}
	j, err := scheduler.NewJob(
		gocron.DurationJob(
			24*time.Hour, //
		),
		gocron.NewTask(
			func() {
				s.removeExpired(context.Background(), zsetKey)
			},
		),
	)
	if err != nil {
		// handle error
	}
	// each job has a unique id
	fmt.Println(j.ID())

	// start the scheduler
	log.Info().Msg("Started job")
	scheduler.Start()
}
