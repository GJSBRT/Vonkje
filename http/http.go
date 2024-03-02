package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	IP 		string `mapstructure:"ip"`
	Port 	int `mapstructure:"port"`
}

type HTTP struct {
	config 	Config
	errChannel 	chan error
	ctx 	context.Context
	log 	*logrus.Logger
	router 	*mux.Router
}

func New(
	config  Config,
	errChannel 	chan error,
	ctx 	context.Context,
	logger 	*logrus.Logger,
) *HTTP {
	return &HTTP{
		config: config,
		errChannel: errChannel,
		ctx: ctx,
		log: logger,
		router: mux.NewRouter(),
	}
}

func (httpServer *HTTP) Start() {
	// Prometheus metrics
	httpServer.router.Handle("/metrics", promhttp.Handler())
	httpServer.router.PathPrefix("/debug/pprof/").HandlerFunc(pprof.Index)
	httpServer.router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	httpServer.router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	httpServer.router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	httpServer.router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Error handlers
	httpServer.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		httpServer.SendErrorResponse(w, "Route not found", http.StatusNotFound)
	})

	httpServer.router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		httpServer.SendErrorResponse(w, "Request method not allowed", http.StatusNotFound)
	})

	server := &http.Server{
		Handler: httpServer.router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func () {
		ipPortCombo := fmt.Sprintf("%s:%d", httpServer.config.IP, httpServer.config.Port)
		httpListener, err := net.Listen("tcp", ipPortCombo)
		if err != nil {
			httpServer.log.WithError(err).Fatal("Failed to start HTTP server")
		}

		httpServer.log.Info("Starting HTTP server on http://" + ipPortCombo)

		server.Serve(httpListener)
	}()

	<-httpServer.ctx.Done()
	httpServer.log.Info("Stopping HTTP servers")

	const timeout = 30 * time.Second
	srvCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	server.Shutdown(srvCtx)
}
