package wave

import (
	"context"

	"github.com/drone/signal"
	"github.com/farwydi/cleanwhale/transport"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// NewWave make Wave use ctx to undo the entire payload chain.
// logger for error notification.
func NewWave(ctx context.Context, logger *zap.Logger) *Wave {
	return &Wave{
		ctx:    signal.WithContext(ctx),
		logger: logger,
	}
}

// Wave workload controller.
type Wave struct {
	ctx    context.Context
	eg     errgroup.Group
	logger *zap.Logger
}

// Add just add payload.
func (s *Wave) Add(f func(ctx context.Context) error) {
	s.eg.Go(func() error {
		return f(s.ctx)
	})
}

// AddIf conditional add payload.
func (s *Wave) AddIf(condition bool, f func(ctx context.Context) error) {
	if condition {
		s.eg.Go(func() error {
			return f(s.ctx)
		})
	}
}

// AddSever just add transport.
func (s *Wave) AddSever(srv transport.Server) {
	s.eg.Go(func() error {
		return srv.Run(s.ctx)
	})
}

// AddSeverIf conditional add transport.
func (s *Wave) AddSeverIf(condition bool, srv transport.Server) {
	if condition {
		s.eg.Go(func() error {
			return srv.Run(s.ctx)
		})
	}
}

// Run waiting all payload in parallel mode
func (s *Wave) Run() {
	if err := s.eg.Wait(); err != nil {
		s.logger.Fatal("program terminated",
			zap.Error(err))
	}
}
