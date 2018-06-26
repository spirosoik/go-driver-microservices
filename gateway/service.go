package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	nsq "github.com/nsqio/go-nsq"
	"github.com/rafaeljesus/nsq-event-bus"
)

type event struct {
	Body map[string]interface{}
}

// Service is a simple  interface for sending driver locations
type Service interface {
	Send(ctx context.Context, vars map[string]string, data []byte, topic string) error
}

type routeService struct {
	emitter *bus.Emitter
	client  *http.Client
}

//NewService factory method
func NewService(e *bus.Emitter, c *http.Client) Service {
	return &routeService{emitter: e, client: c}
}

//Send creates event to be sent in bus
func (s *routeService) Send(_ context.Context, vars map[string]string, data []byte, topic string) error {
	if !nsq.IsValidTopicName(topic) {
		return errors.New("Invalid topic name")
	}
	var body map[string]interface{}
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	for k, v := range vars {
		body[k] = v
	}

	e := event{Body: body}
	if err := s.emitter.Emit(topic, &e); err != nil {
		return err
	}
	return nil
}
