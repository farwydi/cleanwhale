package wave_test

import (
	"context"

	"github.com/farwydi/cleanwhale/config"
	"github.com/farwydi/cleanwhale/log"
	"github.com/farwydi/cleanwhale/transport/whalegrpc"
	"github.com/farwydi/cleanwhale/transport/whalehttp"
	"github.com/farwydi/cleanwhale/wave"
	"google.golang.org/grpc"
)

// This is what your example config looks like
type Config struct {
	Project   config.ProjectConfig
	Transport struct {
		HTTP1 config.HTTPConfig
		HTTP2 config.HTTPConfig
		GRPC  config.GRPCConfig
	}
}

func ExampleWave() {
	var cfg Config
	err := config.LoadConfigs(&cfg, "config1.yml", "config2.yml")
	if err != nil {
		log.Fatal(err)
	}

	logger, err := log.NewLogger(cfg.Project)
	if err != nil {
		log.Fatal(err)
	}

	w := wave.NewWave(context.Background(), logger.Named("main"))

	w.AddSever(whalehttp.NewHTTPServer(
		cfg.Transport.HTTP1, logger.Named("api1"), nil))
	w.AddSever(whalehttp.NewHTTPServer(
		cfg.Transport.HTTP2, logger.Named("api2"), nil))
	w.AddSever(whalegrpc.NewGRPCServer(
		cfg.Transport.GRPC, logger.Named("grpc"), func(srv *grpc.Server) {
			// grpc service is registered here
		}))
	w.AddIf(true, func(ctx context.Context) error {
		<-ctx.Done()
		// any runner here
		return nil
	})

	w.Run()
}
