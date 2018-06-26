package zombie

import (
	"encoding/json"
	"testing"
)

type FakeResponse interface{}

type FakeDriverAPIClient struct {
	GetLocationsByIDFunc func(id string, minutes int, l interface{}) error
}

func (a *FakeDriverAPIClient) GetLocationsByID(id string, minutes int, v interface{}) error {
	if a.GetLocationsByIDFunc != nil {
		return a.GetLocationsByIDFunc(id, minutes, v)
	}
	return nil
}

func TestIsZombie(t *testing.T) {
	api := FakeDriverAPIClient{}
	api.GetLocationsByIDFunc = func(id string, minutes int, v interface{}) error {
		body := `[{"latitude":48.866908,"longitude":2.365218,"updated_at":"2018-06-10T09:06:18Z"},{"latitude":48.866908,"longitude":2.365218,"updated_at":"2018-06-10T09:06:13Z"}]`
		err := json.Unmarshal([]byte(body), v)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	}

	s := NewService(&api)
	z, err := s.IsZombie("1", 5)
	if err != nil {
		t.Fatal(err)
	}

	if z.Status != true {
		t.Errorf("Driver should be zombie but it's not")
	}
}

func TestIsNotZombie(t *testing.T) {
	api := FakeDriverAPIClient{}
	api.GetLocationsByIDFunc = func(id string, minutes int, v interface{}) error {
		body := `[{"latitude":48.866908,"longitude":2.365218,"updated_at":"2018-06-10T09:06:18Z"},{"latitude":48.864271,"longitude":2.350409,"updated_at":"2018-06-10T09:06:14Z"}]`
		err := json.Unmarshal([]byte(body), v)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	}

	s := NewService(&api)
	z, err := s.IsZombie("1", 5)
	if err != nil {
		t.Fatal(err)
	}

	if z.Status != false {
		t.Errorf("Driver should be not zombie but it is")
	}
}
