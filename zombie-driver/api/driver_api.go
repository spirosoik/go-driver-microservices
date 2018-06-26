package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//DriveAPIClient interface for DI
type DriveAPIClient interface {
	GetLocationsByID(id string, minutes int, v interface{}) error
}

//DriverAPI struct which holds HTTP client
type DriverAPI struct {
	client  *http.Client
	baseURL string
}

//New creates a driver API client
func New(c *http.Client, baseURL string) DriveAPIClient {
	api := &DriverAPI{
		client:  c,
		baseURL: baseURL,
	}
	return api
}

//GetLocationsByID returns locations by driver ID
func (api *DriverAPI) GetLocationsByID(id string, minutes int, v interface{}) error {
	endpoint := fmt.Sprintf("%s/drivers/%s/locations?minutes=%d", api.baseURL, id, minutes)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	r, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return nil
}
