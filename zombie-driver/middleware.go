package zombie

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
func (mw loggingMiddleware) IsZombie(id string, minute int) (r *Zombie, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"method": "IsZombie",
			"id":     id,
			"took":   time.Since(begin),
			"data":   err,
		}).Info()
	}(time.Now())
	return mw.next.IsZombie(id, minute)
}
