package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SKF/go-utility/v2/env"
	http_server "github.com/SKF/go-utility/v2/http-server"
	"github.com/SKF/go-utility/v2/log"
	"github.com/gorilla/mux"
	"go.opencensus.io/plugin/ochttp"

	"go-dummy-server/cmd/rest-api/router"
	"go-dummy-server/cmd/rest-api/webserver"
)

var (
	portHealthEndpoint = env.GetAsString("PORT", "9090")
	portWebAPI         = env.GetAsString("PORT", "8666")
)

func main() {

	r := mux.NewRouter()

	router.Router = r

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
		log.WithField("port", portWebAPI).Info("Will start to listen and serve")

		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.WithError(err).Error("HTTP server ListenAndServe")
			sigs <- syscall.SIGTERM
		}
	}()

	go http_server.StartHealthServer(portHealthEndpoint)

	<-serverClosed
	log.Info("Exiting")
}
