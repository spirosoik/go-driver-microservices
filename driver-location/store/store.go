package store

import "github.com/go-redis/redis"

//RedisStore interface for redis client
type RedisStore interface {
	GetByScoreDesc(key string, opt redis.ZRangeBy) ([]string, error)
	AddWithOrder(key string, members ...redis.Z) error
}

//RedisDAO struct which holds redis client
type RedisDAO struct {
	client *redis.Client
}

//NewDAO creates a data access for redis
func NewDAO(addr string) (RedisStore, error) {
	c := redis.NewClient(&redis.Options{Addr: addr})

	_, err := c.Ping().Result()
	if err != nil {
		return nil, err
	}
	redisDAO := &RedisDAO{
		client: c,
	}
	return redisDAO, nil
}

//GetByScoreDesc returns data by key with descending order by score
func (r *RedisDAO) GetByScoreDesc(key string, opt redis.ZRangeBy) ([]string, error) {
	return r.client.ZRevRangeByScore(key, opt).Result()
}

//AddWithOrder adds in specific keys the members in order
func (r *RedisDAO) AddWithOrder(key string, members ...redis.Z) error {
	return r.client.ZAdd(key, members...).Err()
}
