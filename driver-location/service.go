package driver

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/spirosoik/go-driver-microservices/driver-location/store"
)

const storeKey = "locations"

//LocationCreateEvent event from BUS
type locationCreatedEvent struct {
	ID        string  `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type event struct {
	Body locationCreatedEvent
}

//Location response model
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UpdatedAt string  `json:"updated_at"`
}

// Service is a simple  interface for driver locations
type Service interface {
	GetLocations(id string, minute int) (*[]Location, error)
	CreateLocation(event *locationCreatedEvent) error
}

type storeService struct {
	store store.RedisStore
}

//NewService factory method
func NewService(s store.RedisStore) Service {
	return &storeService{store: s}
}

func (s *storeService) GetLocations(id string, minute int) (*[]Location, error) {
	min := time.Now().Add(time.Duration(-minute) * time.Minute).Unix()
	max := time.Now().Unix()

	opt := redis.ZRangeBy{
		Min: fmt.Sprintf("%f", float64(min)),
		Max: fmt.Sprintf("%f", float64(max)),
	}
	r, err := s.store.GetByScoreDesc(createKey(id), opt)
	if err != nil {
		return nil, err
	}
	var buf [][]byte
	for _, v := range r {
		buf = append(buf, []byte(v))
	}

	locations := make([]Location, 0)
	for _, v := range buf {
		var location Location
		err = json.Unmarshal(v, &location)
		if err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}
	return &locations, nil
}

//CreateLocation to store in redis
func (s storeService) CreateLocation(event *locationCreatedEvent) error {
	id := event.ID
	t := time.Now()
	l := Location{
		Latitude:  event.Latitude,
		Longitude: event.Longitude,
		UpdatedAt: t.Format(time.RFC3339),
	}
	value, err := json.Marshal(l)
	if err != nil {
		return err
	}
	score := time.Now().Unix()
	err = s.store.AddWithOrder(createKey(id), redis.Z{
		Score:  float64(score),
		Member: value,
	})
	if err != nil {
		return err
	}
	return nil
}

func createKey(id string) string {
	return fmt.Sprintf("driver:%s", id)
}
