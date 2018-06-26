package main

import (
	"flag"
	"net/http"

	"github.com/sirupsen/logrus"
	zombie "github.com/spirosoik/go-driver-microservices/zombie-driver"
	"github.com/spirosoik/go-driver-microservices/zombie-driver/api"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8082", "HTTP listen address")
	)
	flag.Parse()

	logger := logrus.New()

	//Setup business logic service
	api := api.New(http.DefaultClient, "http://driver-location:8081")
	var s zombie.Service
	{
		s = zombie.NewService(api)
		s = zombie.LoggingMiddleware(logger)(s)
	}

	//Create router
	router := zombie.MakeHTTPHandler(s, logger)

	errchan := make(chan error)
	go func() {
		logger.WithFields(logrus.Fields{
			"protocol": "HTTP",
			"address":  httpAddr,
		}).Info("Set Router Handler")
		errchan <- http.ListenAndServe(*httpAddr, router)
	}()

	logger.WithError(<-errchan).Fatal("Server Error :( !!!!")
}
