package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	rate "github.com/A-Danylevych/btc-api/pkg/ratemicroservice"
	"github.com/A-Danylevych/btc-api/pkg/ratemicroservice/endpoints"
	"github.com/A-Danylevych/btc-api/pkg/ratemicroservice/transport"

	"github.com/go-kit/kit/log"
	"github.com/oklog/oklog/pkg/group"
)

const (
	defaultHTTPPort = "8080"
	url             = "https://api.coinbase.com/v2/prices/spot?currency=UAH"
)

func main() {
	var (
		logger   log.Logger
		httpAddr = net.JoinHostPort("localhost", envString("HTTP_PORT", defaultHTTPPort))
	)

	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	var (
		service     = rate.NewRateService(url)
		eps         = endpoints.NewEndpointSet(service)
		httpHandler = transport.NewHTTPHandler(eps)
	)

	var g group.Group
	{
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
