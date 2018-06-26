package driver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
)

type body interface{}

type scenario struct {
	description  string
	url          string
	method       string
	payload      body
	service      FakeService
	expectedCode int
}

type FakeService struct {
	err      error
	expected []Location
}

func (s FakeService) GetLocations(id string, minute int) (*[]Location, error) {
	if s.err != nil {
		return nil, s.err
	}

	return &s.expected, nil
}

func (s FakeService) CreateLocation(event *locationCreatedEvent) error {
	return nil
}

func TestHTTPHandlerGetLocationsSucces(t *testing.T) {
	tc := scenario{
		description: "Success - Get Driver's locations",
		url:         "/drivers/1/locations?minutes=5",
		method:      "GET",
		service: FakeService{
			expected: []Location{
				{
					Latitude:  48.864193,
					Longitude: 2.350498,
					UpdatedAt: "2018-04-05T22:36:16Z",
				},
				{
					Latitude:  48.863921,
					Longitude: 2.349211,
					UpdatedAt: "2018-04-05T22:36:21Z",
				},
			},
		},
		expectedCode: 200,
	}
	logger := logrus.New()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(tc.method, tc.url, nil)
	if err != nil {
		t.Fatal(err)
	}
	MakeHTTPHandler(&tc.service, logger).ServeHTTP(w, req)

	if status := w.Code; status != tc.expectedCode {
		t.Errorf("returned wrong status code: got %v want %v", status, tc.expectedCode)
	}

	var locations []Location
	err = json.NewDecoder(w.Body).Decode(&locations)
	if err != nil {
		t.Fatal(err)
	}
	expected := tc.service.expected
	if locations == nil {
		t.Errorf("returned unexpected body: got %v want %v", locations, expected)
	}
	if len(locations) != len(expected) {
		t.Error("returned unexpected array sizes")
	}
	for i := range locations {
		if locations[i] != expected[i] {
			t.Errorf("returned unexpected values: got %v want %v", locations[i], expected[i])
		}
	}
}

func TestHTTPHandlerGetLocationsNoEndpoint(t *testing.T) {
	tc := scenario{
		description:  "Fail - Get Driver's locations, missing minutes",
		url:          "/drivers/1/locations",
		method:       "GET",
		service:      FakeService{},
		expectedCode: 404,
	}
	logger := logrus.New()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(tc.method, tc.url, nil)
	if err != nil {
		t.Fatal(err)
	}
	MakeHTTPHandler(&tc.service, logger).ServeHTTP(w, req)

	if status := w.Code; status != tc.expectedCode {
		t.Errorf("returned wrong status code: got %v want %v", status, tc.expectedCode)
	}
}
