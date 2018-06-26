package driver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	bus "github.com/rafaeljesus/nsq-event-bus"
	"github.com/sirupsen/logrus"
)

//backOffTime duration
const backOffTime = time.Duration(10) * time.Minute

//MakeHTTPHandler creates an http handler
func MakeHTTPHandler(s Service, logger *logrus.Logger) http.Handler {
	r := mux.NewRouter()
	r.Path("/drivers/{id}/locations").
		Methods("GET").
		Queries("minutes", "{minutes}").
		HandlerFunc(locationsHandler(s, logger))
	r.Handle("/metrics", prometheus.Handler())
	return r
}

//SetEventHandler to process NSQ events
func SetEventHandler(lookup string, s Service, logger *logrus.Logger) {
	if err := bus.On(bus.ListenerConfig{
		Lookup:  []string{lookup},
		Topic:   "locations",
		Channel: "driver-location",
		HandlerFunc: func(message *bus.Message) (reply interface{}, err error) {
			return process(message, s, logger)
		}}); err != nil {
		logger.WithError(err).Error("Error while consuming message")
	}
}

func process(message *bus.Message, s Service, logger *logrus.Logger) (reply interface{}, err error) {
	str := fmt.Sprintf("%s", message.Payload)
	logger.Info(str)
	e := event{}
	if err = message.DecodePayload(&e); err != nil {
		message.Requeue(backOffTime)
		logger.WithError(err).Error("Decode payload failed. Event Requeued")
		return nil, err
	}
	if err = s.CreateLocation(&e.Body); err != nil {
		logger.WithError(err).Error("Service store failed for Event. Requeued")
		message.Requeue(backOffTime)
		return
	}
	message.Finish()
	return e, nil
}

func locationsHandler(s Service, logger *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		m := r.URL.Query().Get("minutes")

		min, err := strconv.Atoi(m)
		if err != nil {
			logger.WithError(err).Error()
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		d, err := s.GetLocations(vars["id"], min)
		if err != nil {
			logger.WithError(err).Error()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&d); err != nil {
			logger.WithError(err).Error()
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
