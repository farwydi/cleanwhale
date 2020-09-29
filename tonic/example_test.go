package tonic_test

import (
	"github.com/farwydi/cleanwhale/config"
	"github.com/farwydi/cleanwhale/tonic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Config struct {
	Project   config.ProjectConfig
	Transport struct {
		HTTP config.HTTPConfig
	}
}

func ExampleNewMix() {
	var cfg Config

	r := tonic.NewMix(cfg.Project.Mode, zap.L().Named("web"))

	v1 := r.Group("/v1", tonic.V("v1"))
	{
		v1.GET("/hello", func(c *gin.Context) {
			c.Set("handler.url", "/hello")
			// Do something
		})
	}
}
