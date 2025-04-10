package redis_db

import (
	"context"
	"encoding/json"
	"fmt"
	"new-project/services"
	modelServices "new-project/services/entity"
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

func (s *RedisDB) AddCategories(ctx context.Context, cates []modelServices.Categorys) error {

	// tao 2 dataset truoc
	dataMap := make(map[string]string)
	childMap := make(map[string][]string)
	for _, cat := range cates {
		catJson, err := json.Marshal(cat)
		if err != nil {
			return err
		}
		dataMap[cat.CategoryID] = string(catJson)
		parentID := cat.CategoryID
		if cat.Parent.Valid {
			parentID = cat.Parent.Data
		}
		if parentID != cat.CategoryID {
			childMap[parentID] = append(childMap[parentID], cat.CategoryID)
		}
	}
	// Lưu dataMap lên Redis (HMSET)
	if err := s.client.HSet(ctx, CategoryDataKey, dataMap).Err(); err != nil {
		return fmt.Errorf("failed to save category data: %w", err)
	}
	// Lưu childrenMap lên Redis
	for parentID, children := range childMap {
		childrenJson, err := json.Marshal(children)
		if err != nil {
			return fmt.Errorf("failed to marshal children for parent %s: %w", parentID, err)
		}
		if err := s.client.HSet(ctx, CategoryChildrenKey, parentID, childrenJson).Err(); err != nil {
			return fmt.Errorf("failed to save children: %w", err)
		}
	}
	return nil
}
func (s *RedisDB) RemoveCategories(ctx context.Context) error {
	err := s.client.Del(ctx, CategoryChildrenKey).Err()
	if err != nil {
		return err
	}
	return s.client.Del(ctx, CategoryDataKey).Err()

}

// // GetCategoryTree lấy danh sách danh mục dạng cây từ Redis
// func (s *RedisDB) GetCategoryTree(ctx context.Context, rootID string) ([]modelServices.Categorys, error) {
// 	// Lấy toàn bộ dữ liệu danh mục từ Redis
// 	dataMap, err := s.client.HGetAll(ctx, CategoryDataKey).Result()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get category data: %w", err)
// 	}

// 	// Lấy toàn bộ dữ liệu con từ Redis
// 	childMap, err := s.client.HGetAll(ctx, CategoryChildrenKey).Result()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get category children: %w", err)
// 	}
// 	fmt.Println("len childmap", len(childMap))
// 	// Parse dữ liệu thành map các danh mục
// 	categories := make(map[string]*modelServices.Categorys)
// 	for _, jsonStr := range dataMap {
// 		var cat modelServices.Categorys
// 		if err := json.Unmarshal([]byte(jsonStr), &cat); err != nil {
// 			return nil, fmt.Errorf("failed to unmarshal category: %w", err)
// 		}
// 		categories[cat.CategoryID] = &cat
// 	}

// 	// Xây dựng cây bằng cách gắn các con vào cha
// 	for parentID, childrenJson := range childMap {
// 		var childrenIDs []string
// 		if err := json.Unmarshal([]byte(childrenJson), &childrenIDs); err != nil {
// 			return nil, fmt.Errorf("failed to unmarshal children for parent %s: %w", parentID, err)
// 		}

// 		if parent, exists := categories[parentID]; exists {
// 			var childCategories []modelServices.Categorys
// 			for _, childID := range childrenIDs {
// 				if child, exists := categories[childID]; exists {
// 					childCategories = append(childCategories, *child)
// 				}
// 			}
// 			parent.Childs = modelServices.Narg[[]modelServices.Categorys]{Data: childCategories, Valid: true}
// 		}
// 	}

// 	// Nếu có rootID, chỉ lấy danh mục đó và các con của nó
// 	if rootID != "" {
// 		if rootCat, exists := categories[rootID]; exists {
// 			return []modelServices.Categorys{*rootCat}, nil
// 		}
// 		return nil, fmt.Errorf("category with ID %s not found", rootID)
// 	}

// 	// Nếu không có rootID, lấy toàn bộ danh mục gốc (không có parent hoặc parent không tồn tại)
// 	var result []modelServices.Categorys
// 	for _, cat := range categories {
// 		if !cat.Parent.Valid || categories[cat.Parent.Data] == nil {
// 			result = append(result, *cat)
// 		}
// 	}

// 	return result, nil
// }

func (s *RedisDB) GetCategoryTree(ctx context.Context, rootID string) ([]modelServices.Categorys, error) {
	// Lấy toàn bộ dữ liệu danh mục từ Redis
	dataMap, err := s.client.HGetAll(ctx, CategoryDataKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get category data: %w", err)
	}

	// Lấy toàn bộ dữ liệu con từ Redis
	childMap, err := s.client.HGetAll(ctx, CategoryChildrenKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get category children: %w", err)
	}

	// Parse dữ liệu thành map các danh mục
	categories := make(map[string]*modelServices.Categorys)
	for _, jsonStr := range dataMap {
		var cat modelServices.Categorys
		if err := json.Unmarshal([]byte(jsonStr), &cat); err != nil {
			return nil, fmt.Errorf("failed to unmarshal category: %w", err)
		}
		categories[cat.CategoryID] = &cat
	}

	// Hàm đệ quy để xây dựng cây con với độ sâu tối đa 3 lớp
	var buildTree func(parentID string, currentDepth int) ([]modelServices.Categorys, error)
	buildTree = func(parentID string, currentDepth int) ([]modelServices.Categorys, error) {
		childrenJson, exists := childMap[parentID]
		if !exists || currentDepth >= 3 { // Giới hạn độ sâu 3 lớp
			return nil, nil
		}

		var childrenIDs []string
		if err := json.Unmarshal([]byte(childrenJson), &childrenIDs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal children for parent %s: %w", parentID, err)
		}

		var childCategories []modelServices.Categorys
		for _, childID := range childrenIDs {
			if child, exists := categories[childID]; exists {
				// Tạo bản sao của category để tránh tham chiếu trực tiếp
				childCopy := *child

				// Đệ quy xây dựng cây con với độ sâu tăng lên
				grandchildren, err := buildTree(childID, currentDepth+1)
				if err != nil {
					return nil, err
				}

				if len(grandchildren) > 0 {
					childCopy.Childs = modelServices.Narg[[]modelServices.Categorys]{
						Data:  grandchildren,
						Valid: true,
					}
				}

				childCategories = append(childCategories, childCopy)
			}
		}

		return childCategories, nil
	}

	// Nếu có rootID, chỉ lấy danh mục đó và các con của nó (tối đa 3 lớp)
	if rootID != "" {
		rootCat, exists := categories[rootID]
		if !exists {
			return nil, fmt.Errorf("category with ID %s not found", rootID)
		}

		// Tạo bản sao của root category
		rootCopy := *rootCat

		// Xây dựng cây con với độ sâu bắt đầu từ 1 (vì root là lớp 0)
		children, err := buildTree(rootID, 1)
		if err != nil {
			return nil, err
		}

		if len(children) > 0 {
			rootCopy.Childs = modelServices.Narg[[]modelServices.Categorys]{
				Data:  children,
				Valid: true,
			}
		}

		return []modelServices.Categorys{rootCopy}, nil
	}

	// Nếu không có rootID, lấy toàn bộ danh mục gốc và xây dựng cây con
	var result []modelServices.Categorys
	for _, cat := range categories {
		if !cat.Parent.Valid || categories[cat.Parent.Data] == nil {
			// Tạo bản sao của category gốc
			rootCopy := *cat

			// Xây dựng cây con với độ sâu bắt đầu từ 1
			children, err := buildTree(cat.CategoryID, 1)
			if err != nil {
				return nil, err
			}

			if len(children) > 0 {
				rootCopy.Childs = modelServices.Narg[[]modelServices.Categorys]{
					Data:  children,
					Valid: true,
				}
			}

			result = append(result, rootCopy)
		}
	}

	return result, nil
}
