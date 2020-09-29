package tonic

import (
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	shMx        sync.RWMutex
	skipHandler = map[string]struct{}{
		healthPath:                           {},
		metricsPath:                          {},
		defaultPProfPrefix + "/cmdline":      {},
		defaultPProfPrefix + "/profile":      {},
		defaultPProfPrefix + "/symbol":       {},
		defaultPProfPrefix + "/trace":        {},
		defaultPProfPrefix + "/allocs":       {},
		defaultPProfPrefix + "/block":        {},
		defaultPProfPrefix + "/goroutine":    {},
		defaultPProfPrefix + "/heap":         {},
		defaultPProfPrefix + "/mutex":        {},
		defaultPProfPrefix + "/threadcreate": {},
	}

	// V alias declares handler.version with the given version.
	V = func(version string) func(c *gin.Context) {
		return func(c *gin.Context) {
			c.Set("handler.version", version)
		}
	}
)

func inSkipHandler(path string) bool {
	shMx.RLock()
	defer shMx.RUnlock()
	_, found := skipHandler[path]
	return found
}

// RegisterSkipHandler adds skipPath ignore paths to handlers like logger and recover.
func RegisterSkipHandler(skipPath []string) {
	shMx.Lock()
	defer shMx.Unlock()
	for _, path := range skipPath {
		skipHandler[path] = struct{}{}
	}
}
