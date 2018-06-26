package driver

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/go-redis/redis"
)

type FakeRedisStore struct {
	GetByScoreDescFunc func(key string, opt redis.ZRangeBy) ([]string, error)
	AddWithOrderFunc   func(key string, members ...redis.Z) error
}

func (r *FakeRedisStore) GetByScoreDesc(key string, opt redis.ZRangeBy) ([]string, error) {
	if r.GetByScoreDescFunc != nil {
		return r.GetByScoreDescFunc(key, opt)
	}
	return nil, nil
}

func (r *FakeRedisStore) AddWithOrder(key string, members ...redis.Z) error {
	if r.AddWithOrderFunc != nil {
		return r.AddWithOrderFunc(key, members...)
	}
	return nil
}
func TestGetLocationsListSuccess(t *testing.T) {
	store := FakeRedisStore{
		GetByScoreDescFunc: func(key string, opt redis.ZRangeBy) ([]string, error) {
			l1 := Location{
				Latitude:  48.864193,
				Longitude: 2.350498,
				UpdatedAt: "2018-04-05T22:36:17Z",
			}
			l2 := Location{
				Latitude:  48.864193,
				Longitude: 2.350498,
				UpdatedAt: "2018-04-05T22:36:21Z",
			}
			out1, err := json.Marshal(l1)
			if err != nil {
				t.Fatal(err)
			}
			out2, err := json.Marshal(l2)
			if err != nil {
				t.Fatal(err)
			}
			return []string{string(out1), string(out2)}, nil
		},
	}

	s := NewService(&store)
	l, err := s.GetLocations("1", 10)
	if err != nil {
		t.Fatal(err)
	}

	if len(*l) != 2 {
		t.Errorf("Locations are less than expected: got %d expected 2", len(*l))
	}
}

func TestGetLocationsError(t *testing.T) {
	store := FakeRedisStore{
		GetByScoreDescFunc: func(key string, opt redis.ZRangeBy) ([]string, error) {
			return nil, errors.New("Error fetching locations from Redis")
		},
	}
	s := NewService(&store)
	_, err := s.GetLocations("1", 10)
	if err == nil {
		t.Fatal(err)
	}
}
func TestCreateLocationSuccess(t *testing.T) {
	store := FakeRedisStore{
		AddWithOrderFunc: func(key string, members ...redis.Z) error {
			return nil
		},
	}
	s := NewService(&store)
	err := s.CreateLocation(&locationCreatedEvent{
		Latitude:  48.864193,
		Longitude: 2.350498,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateCreateLocationError(t *testing.T) {
	store := FakeRedisStore{
		AddWithOrderFunc: func(key string, members ...redis.Z) error {
			return errors.New("Failed to store in Redis")
		},
	}
	s := NewService(&store)
	err := s.CreateLocation(&locationCreatedEvent{
		Latitude:  48.864193,
		Longitude: 2.350498,
	})
	if err == nil {
		t.Fatal(err)
	}
}
