package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opencensus.io/plugin/ochttp"

	"go-dummy-server/cmd/rest-api/logger"
	"go-dummy-server/cmd/rest-api/router"
	"go-dummy-server/cmd/rest-api/webserver"
)

var (
	portWebAPI = os.Getenv("PORT")
)

func main() {
	if portWebAPI == "" {
		portWebAPI = "8666"
	}

	log := logger.Log()

	r := router.Get()

	webserver.NewEndpoint("AddEndpoint", r).
		Path("/addEndpoint").
		Methods(http.MethodPost).
		HandlerFunc(webserver.AddEndpoint)

	httpServer := &http.Server{
		Addr:           ":" + portWebAPI,
		Handler:        &ochttp.Handler{Handler: r},
		ReadTimeout:    30 * time.Second, //nolint:gomnd
		WriteTimeout:   30 * time.Second, //nolint:gomnd
		MaxHeaderBytes: 1 << 20,          //nolint:gomnd
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	serverClosed := make(chan struct{}, 1)

	go func() {
		<-sigs
		log.Info("Will try to exit gracefully")
		close(serverClosed)
	}()

	go func() {
		log.Infow("Will start to listen and serve",
			"port", portWebAPI)

		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Errorw("HTTP server ListenAndServe",
				"error", err)
			sigs <- syscall.SIGTERM
		}
	}()

	<-serverClosed
	log.Info("Exiting")
}
