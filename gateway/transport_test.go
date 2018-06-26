package gateway

import (
	"bytes"
	"context"
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
	expected *http.Response
}

func (s FakeService) Send(_ context.Context, vars map[string]string, data []byte, topic string) error {
	if s.err != nil {
		return s.err
	}

	return nil
}

func TestHTTPHandlerGetDriverWrongMethod(t *testing.T) {
	u1 := URL{
		Path:   "/drivers/{id}",
		Method: "GET",
	}
	u1.HTTP.Host = "zombie-driver"
	u1.HTTP.Port = "8082"

	urls := []URL{u1}

	tc := scenario{
		description:  "Fail - Get Driver wrong method",
		url:          "/drivers/1",
		method:       "PATCH",
		service:      FakeService{},
		expectedCode: 405,
	}
	logger := logrus.New()
	ctx := context.Background()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(tc.method, tc.url, nil)
	if err != nil {
		t.Fatal(err)
	}
	MaketHTTPHandler(ctx, urls, &tc.service, logger).ServeHTTP(w, req)

	if status := w.Code; status != tc.expectedCode {
		t.Errorf("returned wrong status code: got %v want %v", status, tc.expectedCode)
	}
}

func TestHTTPHandlerUpdateLocation(t *testing.T) {
	u1 := URL{
		Path:   "/drivers/{id}/locations",
		Method: "PATCH",
	}
	u1.HTTP.Host = ""
	u1.NSQ.Topic = "locations"

	urls := []URL{u1}

	tests := []scenario{
		{
			description: "Success - Update Driver location",
			url:         "/drivers/1/locations",
			method:      "PATCH",
			service:     FakeService{},
			payload: event{
				Body: map[string]interface{}{
					"latitude": 48.864193, "longitude": 2.350498,
				},
			},
			expectedCode: 200,
		},
		{
			description: "Fail - Update Driver location, non payload",
			url:         "/drivers/1/locations",
			method:      "POST",
			service:     FakeService{},
			payload: event{
				Body: map[string]interface{}{
					"Latitude": 48.864193, "Longitude": 2.350498,
				},
			},
			expectedCode: 405,
		},
	}
	logger := logrus.New()
	ctx := context.Background()
	w := httptest.NewRecorder()
	for _, tc := range tests {
		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(tc.payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest(tc.method, tc.url, buf)
		if err != nil {
			t.Fatal(err)
		}
		MaketHTTPHandler(ctx, urls, &tc.service, logger).ServeHTTP(w, req)

		if status := w.Code; status != tc.expectedCode {
			t.Errorf("returned wrong status code: got %v want %v", status, tc.expectedCode)
		}
	}
}
