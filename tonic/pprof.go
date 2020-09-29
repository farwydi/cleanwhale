package tonic

import (
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

const (
	// defaultPProfPrefix url prefix of pprof
	defaultPProfPrefix = "/__/pprof"
)

// AddPProfHandler the standard HandlerFuncs from the net/http/pprof package with
// the provided gin.Group with defaultPProfPrefix.
func AddPProfHandler(r *gin.Engine) {
	prefixRouter := r.Group(defaultPProfPrefix)
	{
		prefixRouter.GET("/", gin.WrapF(pprof.Index))
		prefixRouter.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		prefixRouter.GET("/profile", gin.WrapF(pprof.Profile))
		prefixRouter.POST("/symbol", gin.WrapF(pprof.Symbol))
		prefixRouter.GET("/symbol", gin.WrapF(pprof.Symbol))
		prefixRouter.GET("/trace", gin.WrapF(pprof.Trace))
		prefixRouter.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		prefixRouter.GET("/block", gin.WrapH(pprof.Handler("block")))
		prefixRouter.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		prefixRouter.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		prefixRouter.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		prefixRouter.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}
}
