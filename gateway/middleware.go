package gateway

import (
	"context"
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

func (mw loggingMiddleware) Send(ctx context.Context, vars map[string]string, data []byte, topic string) error {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"method": "Send",
			"id":     vars,
			"took":   time.Since(begin),
			"topic":  topic,
		}).Info()
	}(time.Now())
	return mw.next.Send(ctx, vars, data, topic)
}
