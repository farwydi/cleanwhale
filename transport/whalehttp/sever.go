// Package whalehttp defines the http implementation of the transport.Server transport.
package whalehttp

import (
	"context"
	"net/http"
	"time"

	"github.com/farwydi/cleanwhale/config"
	"github.com/farwydi/cleanwhale/transport"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// NewHTTPServer Creates an http transport based on router.
// Configures the main Addr parameter with cfg.
// logger is used for startup notification and shutting.
func NewHTTPServer(cfg config.HTTPConfig, logger *zap.Logger, router http.Handler) transport.Server {
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	return &httpServer{srv: srv, logger: logger, addr: cfg.Addr}
}

type httpServer struct {
	addr   string
	srv    *http.Server
	logger *zap.Logger
}

// Run ListenAndServe http server and also expects ctx.Done for graceful shutdown.
// Shutdown is limited to 5s timeout.
func (t *httpServer) Run(ctx context.Context) error {
	var g errgroup.Group
	g.Go(func() error {
		<-ctx.Done()

		t.logger.Info("Shutting down http server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return t.srv.Shutdown(ctx)
	})
	g.Go(func() error {
		t.logger.Info("Starting http server",
			zap.String("addr", t.addr))
		if err := t.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	return g.Wait()
}
