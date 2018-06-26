package zombie

import (
	"encoding/json"
	"errors"
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
	expected Zombie
}

func (s FakeService) IsZombie(id string, minute int) (*Zombie, error) {
	if s.err != nil {
		return nil, s.err
	}

	return &s.expected, nil
}

func TestHTTPHandlerIsZombieInternalError(t *testing.T) {
	tc := scenario{
		description: "Fail - Get Driver's state failed",
		url:         "/drivers/1",
		method:      "GET",
		service: FakeService{
			err: errors.New("Failed to error"),
		},
		expectedCode: 500,
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

func TestHTTPHandlerIsZombieSuccess(t *testing.T) {
	tc := scenario{
		description: "Success - Get Driver's state is not zombie",
		url:         "/drivers/1",
		method:      "GET",
		service: FakeService{
			expected: Zombie{
				ID:     1,
				Status: false,
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

	var z Zombie
	err = json.NewDecoder(w.Body).Decode(&z)
	if err != nil {
		t.Fatal(err)
	}
	expected := tc.service.expected
	if &z == nil {
		t.Errorf("returned unexpected body: got %v want %v", z, expected)
	}
	if z != expected {
		t.Errorf("returned unexpected values: got %v want %v", z, expected)
	}
}
