package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	bus "github.com/rafaeljesus/nsq-event-bus"
	"github.com/sirupsen/logrus"
	"github.com/spirosoik/go-driver-microservices/gateway"
)

//Config for gateway
type Config struct {
	Urls []gateway.URL `config:"urls"`
}

func main() {
	var (
		config   = flag.String("config", "config.yaml", "Yaml config file")
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
		nsqAddr  = flag.String("nsq.addr", ":4150", "NSQ listen address")
	)
	flag.Parse()

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Load configuration
	cfg := Config{}
	loader := confita.NewLoader(
		file.NewBackend(*config),
	)
	err := loader.Load(context.Background(), &cfg)
	if err != nil {
		logger.WithError(err).Fatal("failed to load configuration")
	}

	// Emitter for bus setup
	emitter, err := bus.NewEmitter(bus.EmitterConfig{
		Address:     *nsqAddr,
		MaxInFlight: 100,
	})
	if err != nil {
		logger.WithError(err).Fatal("failed to create BUS emitter")
	}

	//Set service
	var s gateway.Service
	{
		s = gateway.NewService(emitter, http.DefaultClient)
		s = gateway.LoggingMiddleware(logger)(s)
	}

	// Handler setup
	ctx := context.Background()
	h := gateway.MaketHTTPHandler(ctx, cfg.Urls, s, logger)

	// Run gateway
	errchan := make(chan error)
	go func() {
		logger.WithFields(logrus.Fields{
			"protocol": "HTTP",
			"address":  httpAddr,
		}).Info("Set Router Handler")
		errchan <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Error("Server Error :( !!!!", <-errchan)
}
