package driver

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Middleware describes a service middleware
type Middleware func(Service) Service

//LoggingMiddleware create a new logging middleware
func LoggingMiddleware(logger *logrus.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger *logrus.Logger
}

//GetLocations middleware logging for service
func (mw loggingMiddleware) GetLocations(id string, minute int) (r *[]Location, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"method": "GetLocations",
			"id":     id,
			"took":   time.Since(begin),
			"data":   err,
		}).Info()
	}(time.Now())
	return mw.next.GetLocations(id, minute)
}

//GetLocations middleware logging for service
func (mw loggingMiddleware) CreateLocation(event *locationCreatedEvent) (err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"method": "CreateLocation",
			"event":  event,
			"took":   time.Since(begin),
			"err":    err,
		}).Info()
	}(time.Now())
	return mw.next.CreateLocation(event)
}
