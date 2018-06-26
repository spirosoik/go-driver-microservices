package main

import (
	"flag"
	"net/http"

	"github.com/spirosoik/go-driver-microservices/driver-location/store"

	"github.com/sirupsen/logrus"

	"github.com/spirosoik/go-driver-microservices/driver-location"
)

func main() {
	var (
		httpAddr   = flag.String("http.addr", ":8081", "HTTP listen address")
		lookupAddr = flag.String("lookup.addr", ":4161", "NSQ lookup address")
		redisAddr  = flag.String("redis.addr", ":6379", "Redis address")
	)
	flag.Parse()

	logger := logrus.New()

	if *lookupAddr == "" {
		logger.Fatal("Lookup address is needed")
	}

	if *redisAddr == "" {
		logger.Fatal("Redis address is needed")
	}

	//Setup business logic service
	r, err := store.NewDAO(*redisAddr)
	if err != nil {
		logger.WithError(err).Fatal("Redis failed to created")
	}
	var s driver.Service
	{
		s = driver.NewService(r)
		s = driver.LoggingMiddleware(logger)(s)
	}
	//set event handler
	driver.SetEventHandler(*lookupAddr, s, logger)

	//Create router
	router := driver.MakeHTTPHandler(s, logger)

	errchan := make(chan error)
	go func() {
		logger.WithFields(logrus.Fields{
			"protocol": "HTTP",
			"address":  httpAddr,
		}).Info("Set Router Handler")
		errchan <- http.ListenAndServe(*httpAddr, router)
	}()

	logger.Error("Server Error :( !!!!", <-errchan)
}
