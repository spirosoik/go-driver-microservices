package zombie

import (
	"strconv"

	geo "github.com/kellydunn/golang-geo"
	"github.com/spirosoik/go-driver-microservices/zombie-driver/api"
)

//Location DTO
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UpdatedAt string  `json:"updated_at"`
}

//Zombie DTO
type Zombie struct {
	ID     int  `json:"id"`
	Status bool `json:"zombie"`
}

// Service is a simple  interface for driver locations
type Service interface {
	IsZombie(id string, minute int) (*Zombie, error)
}

type httpService struct {
	client api.DriveAPIClient
}

//NewService factory method
func NewService(api api.DriveAPIClient) Service {
	return &httpService{client: api}
}

//IsZombie checks if a driver is zombie or not
func (s *httpService) IsZombie(id string, minute int) (*Zombie, error) {
	var ls []Location
	err := s.client.GetLocationsByID(id, minute, &ls)
	if err != nil {
		return nil, err
	}
	var d float64
	var point *geo.Point
	for _, l := range ls {
		p := geo.NewPoint(l.Latitude, l.Longitude)
		if point != nil {
			d += point.GreatCircleDistance(p)
		}
		point = p
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	meters := d * 1000
	if meters < 500 {
		return &Zombie{ID: idInt, Status: true}, nil
	}
	return &Zombie{ID: idInt, Status: false}, nil
}
